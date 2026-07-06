#!/usr/bin/env python3
"""
Mobile client smoke test for the Liar's Bar server.

It verifies the gateway path mobile apps should use:
  - HTTP API through Nginx: /api/v1
  - WebSocket through Nginx: /ws?token=...

Usage:
  python3 scripts/mobile_smoke_test.py
  python3 scripts/mobile_smoke_test.py http://SERVER_IP:8081
"""

from __future__ import annotations

import base64
import hashlib
import json
import os
import random
import socket
import ssl
import struct
import sys
import time
import urllib.error
import urllib.parse
import urllib.request


BASE_URL = sys.argv[1].rstrip("/") if len(sys.argv) > 1 else "http://localhost:8081"
API_BASE = f"{BASE_URL}/api/v1"


def request(method: str, path: str, body: dict | None = None, token: str | None = None) -> dict:
    data = None if body is None else json.dumps(body).encode("utf-8")
    req = urllib.request.Request(f"{API_BASE}{path}", data=data, method=method)
    req.add_header("Accept", "application/json")
    if body is not None:
        req.add_header("Content-Type", "application/json")
    if token:
        req.add_header("Authorization", f"Bearer {token}")

    try:
        with urllib.request.urlopen(req, timeout=8) as resp:
            raw = resp.read().decode("utf-8")
            return json.loads(raw) if raw else {}
    except urllib.error.HTTPError as exc:
        detail = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"{method} {path} failed: HTTP {exc.code} {detail}") from exc


def ws_handshake(token: str) -> None:
    parsed = urllib.parse.urlparse(BASE_URL)
    secure = parsed.scheme == "https"
    host = parsed.hostname or "localhost"
    port = parsed.port or (443 if secure else 80)
    query_token = urllib.parse.quote(token, safe="")
    path = f"/ws?token={query_token}"
    key = base64.b64encode(os.urandom(16)).decode("ascii")

    raw_sock = socket.create_connection((host, port), timeout=8)
    sock = ssl.create_default_context().wrap_socket(raw_sock, server_hostname=host) if secure else raw_sock
    try:
        request_text = (
            f"GET {path} HTTP/1.1\r\n"
            f"Host: {host}:{port}\r\n"
            "Upgrade: websocket\r\n"
            "Connection: Upgrade\r\n"
            f"Sec-WebSocket-Key: {key}\r\n"
            "Sec-WebSocket-Version: 13\r\n"
            "\r\n"
        )
        sock.sendall(request_text.encode("ascii"))
        response = sock.recv(4096).decode("iso-8859-1", errors="replace")
        if " 101 " not in response.split("\r\n", 1)[0]:
            raise RuntimeError(f"WebSocket upgrade failed: {response.splitlines()[0] if response else 'empty response'}")

        expected = base64.b64encode(
            hashlib.sha1((key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11").encode("ascii")).digest()
        ).decode("ascii")
        if expected not in response:
            raise RuntimeError("WebSocket accept key mismatch")
    finally:
        try:
            sock.close()
        except OSError:
            pass


def main() -> None:
    suffix = f"{int(time.time())}{random.randint(1000, 9999)}"
    username = f"mobile_{suffix}"
    password = "MobileTest123"
    nickname = f"Mobile Test {suffix[-4:]}"

    print(f"base_url={BASE_URL}")

    register = request("POST", "/auth/register", {
        "username": username,
        "password": password,
        "nickname": nickname,
    })
    assert register.get("code") == 0, register
    print("register=ok")

    login = request("POST", "/auth/login", {"username": username, "password": password})
    token = login.get("token")
    assert login.get("code") == 0 and token, login
    print("login=ok")

    profile = request("GET", "/user/profile", token=token)
    assert profile.get("code") == 0, profile
    print("profile=ok")

    lobby = request("GET", "/lobby", token=token)
    assert lobby.get("code") == 0, lobby
    print("lobby=ok")

    room = request("POST", "/rooms", {"name": f"Mobile Smoke {suffix[-4:]}"}, token=token)
    room_id = room.get("data", {}).get("ID") or room.get("data", {}).get("id")
    assert room.get("code") == 0 and room_id, room
    print(f"create_room=ok room_id={room_id}")

    join = request("POST", f"/rooms/{room_id}/join", token=token)
    assert join.get("code") == 0, join
    print("join_room=ok")

    ws_handshake(token)
    print("websocket_upgrade=ok")

    leave = request("POST", f"/rooms/{room_id}/leave", token=token)
    assert leave.get("code") == 0, leave
    print("leave_room=ok")

    print("mobile_smoke_test=passed")


if __name__ == "__main__":
    main()
