#!/usr/bin/env bash
# ============================================================
# Liar's Bar 骗子酒馆 - 一键运行脚本
# 基于 docker-compose 编排所有服务(MySQL/Redis/AI/Backend/Frontend/Nginx)
#
# 用法:
#   ./run.sh          # 启动(默认)
#   ./run.sh start    # 启动
#   ./run.sh stop     # 停止
#   ./run.sh restart  # 重启
#   ./run.sh status   # 查看状态
#   ./run.sh logs     # 查看所有服务日志(实时)
#   ./run.sh logs <服务名>  # 查看指定服务日志,如 ./run.sh logs backend
#   ./run.sh clean    # 停止并删除容器/网络(保留数据卷)
# ============================================================
set -e

# ---------- 颜色输出 ----------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

info()    { echo -e "${BLUE}[INFO]${NC}  $*"; }
ok()      { echo -e "${GREEN}[OK]${NC}    $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}  $*"; }
err()     { echo -e "${RED}[ERROR]${NC} $*"; }
title()   { echo -e "\n${CYAN}==== $* ====${NC}"; }

# ---------- 路径与常量 ----------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

COMPOSE_FILE="docker-compose.yml"
PROJECT_NAME="liarsbar"

# 对外端口(与 docker-compose.yml 保持一致)
NGINX_PORT=8081      # 总入口(推荐访问)
FRONTEND_PORT=3000
BACKEND_PORT=8082
AI_PORT=8000
MYSQL_PORT=3307
REDIS_PORT=6379

# ---------- 前置检查 ----------
check_prerequisites() {
    title "前置检查"

    if ! command -v docker >/dev/null 2>&1; then
        err "未找到 docker,请先安装 Docker"
        exit 1
    fi
    if ! command -v docker-compose >/dev/null 2>&1 && ! docker compose version >/dev/null 2>&1; then
        err "未找到 docker-compose,请先安装 Docker Compose"
        exit 1
    fi

    if ! docker info >/dev/null 2>&1; then
        err "Docker daemon 未运行,请先启动 Docker 服务(sudo systemctl start docker)"
        exit 1
    fi
    ok "Docker 运行正常"
}

# 统一调用 docker compose(优先 v2,回退 v1)
dc() {
    if docker compose version >/dev/null 2>&1; then
        docker compose -p "$PROJECT_NAME" "$@"
    else
        docker-compose -p "$PROJECT_NAME" "$@"
    fi
}

# ---------- 端口占用检查 ----------
check_ports() {
    local busy=()
    for p in "$NGINX_PORT" "$FRONTEND_PORT" "$BACKEND_PORT" "$AI_PORT" "$MYSQL_PORT" "$REDIS_PORT"; do
        if ss -ltn 2>/dev/null | grep -q ":${p}\b" || netstat -ltn 2>/dev/null | grep -q ":${p}\b"; then
            busy+=("$p")
        fi
    done
    if [ ${#busy[@]} -gt 0 ]; then
        warn "以下端口可能已被占用: ${busy[*]}"
        warn "若为同项目的旧容器,脚本会先停止它们;若为其他程序,请手动处理后再运行"
    fi
}

# ---------- 启动 ----------
start() {
    check_prerequisites
    check_ports

    title "构建并启动容器(首次会比较慢,需要拉镜像/编译)"
    dc -f "$COMPOSE_FILE" up -d --build

    title "等待服务就绪"
    wait_for_nginx

    title "服务状态"
    dc -f "$COMPOSE_FILE" ps

    print_access_info
}

# ---------- 等待 Nginx 入口可用 ----------
wait_for_nginx() {
    info "等待 Nginx 入口 http://localhost:${NGINX_PORT} 响应..."
    local i
    for i in $(seq 1 60); do
        if curl -fsS "http://localhost:${NGINX_PORT}/" >/dev/null 2>&1; then
            ok "Nginx 已就绪(第 ${i} 次尝试)"
            return 0
        fi
        printf "."
        sleep 2
    done
    echo
    warn "Nginx 在 120s 内未响应,可能后端仍在启动,请用 './run.sh status' 查看状态"
    warn "也可以用 './run.sh logs' 排查具体服务日志"
}

# ---------- 打印访问信息 ----------
print_access_info() {
    title "访问地址"
    echo -e "  ${GREEN}前端页面(推荐) :${NC} http://localhost:${NGINX_PORT}"
    echo -e "  ${GREEN}前端直连       :${NC} http://localhost:${FRONTEND_PORT}"
    echo -e "  ${GREEN}移动端统一入口 :${NC} http://服务器IP:${NGINX_PORT}"
    echo -e "  后端 API(经 Nginx)  : http://服务器IP:${NGINX_PORT}/api/v1"
    echo -e "  WebSocket(经 Nginx) : ws://服务器IP:${NGINX_PORT}/ws?token=JWT"
    echo -e "  后端 API(直连调试)  : http://localhost:${BACKEND_PORT}/api/v1"
    echo -e "  AI 服务             : http://localhost:${AI_PORT}/docs"
    echo
    echo -e "  ${CYAN}首次使用需要在页面注册账号后登录${NC}"
    echo
    echo -e "  常用命令:"
    echo -e "    ./run.sh status    查看状态"
    echo -e "    ./run.sh logs      查看日志"
    echo -e "    ./run.sh stop      停止"
    echo -e "    ./run.sh restart   重启"
}

# ---------- 停止 ----------
stop() {
    title "停止服务"
    dc -f "$COMPOSE_FILE" down
    ok "已停止"
}

# ---------- 重启 ----------
restart() {
    stop
    start
}

# ---------- 状态 ----------
status() {
    title "服务状态"
    dc -f "$COMPOSE_FILE" ps
    echo
    title "入口连通性"
    if curl -fsS "http://localhost:${NGINX_PORT}/" >/dev/null 2>&1; then
        ok "Nginx 入口 http://localhost:${NGINX_PORT} 可访问"
    else
        err "Nginx 入口 http://localhost:${NGINX_PORT} 不可访问"
    fi
    if curl -fsS "http://localhost:${NGINX_PORT}/health" >/dev/null 2>&1; then
        ok "移动端网关健康检查 http://localhost:${NGINX_PORT}/health 可访问"
    else
        err "移动端网关健康检查失败"
    fi
    if curl -fsS "http://localhost:${BACKEND_PORT}/api/v1/lobby" >/dev/null 2>&1; then
        ok "后端 API 直连可访问"
    else
        warn "后端 API 探测失败(可能需要登录 token,仅供参考)"
    fi
}

# ---------- 日志 ----------
logs() {
    if [ -n "$1" ]; then
        dc -f "$COMPOSE_FILE" logs -f --tail=200 "$1"
    else
        dc -f "$COMPOSE_FILE" logs -f --tail=200
    fi
}

# ---------- 清理(保留数据卷) ----------
clean() {
    title "停止并删除容器/网络(保留数据卷)"
    dc -f "$COMPOSE_FILE" down
    ok "已清理(数据卷 mysql_data 保留)"
    warn "如需彻底删除数据,执行: docker volume rm ${PROJECT_NAME}_mysql_data"
}

# ---------- 帮助 ----------
usage() {
    cat <<EOF
Liar's Bar 骗子酒馆 - 一键运行脚本

用法:
  ./run.sh [命令]

命令:
  start     启动所有服务(默认)
  stop      停止所有服务
  restart   重启所有服务
  status    查看运行状态与入口连通性
  logs [服务名]  查看日志(实时),可指定服务如 backend/frontend/mysql/redis/ai-service/nginx
  clean     停止并删除容器/网络(保留数据卷)
  help      显示本帮助

无参数等同于 start。
EOF
}

# ---------- 入口 ----------
case "${1:-start}" in
    start)   start ;;
    stop)    stop ;;
    restart) restart ;;
    status)  status ;;
    logs)    shift; logs "$@" ;;
    clean)   clean ;;
    help|-h|--help) usage ;;
    *)
        err "未知命令: $1"
        usage
        exit 1
        ;;
esac
