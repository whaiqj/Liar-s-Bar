再写

# Liar's Bar Online · 安卓客户端 UI 设计提示词

本文件面向安卓客户端的 UI 设计，基于现有 Web 端（`frontend/src`）的视觉语言提炼，供 Midjourney / Stable Diffusion / Figma AI / Sora-style 图像模型直接使用。

---

## 0. 设计系统摘要（Style Tokens）

> 所有分屏提示词都默认继承以下风格，确保整套 UI 视觉一致。

| 维度       | 取值                                                                                           | 含义                  |
| ---------- | ---------------------------------------------------------------------------------------------- | --------------------- |
| 主题       | 复古酒馆 + 暗黑赌局 + 心理博弈                                                                 | Liar's Bar / 骗子酒馆 |
| 形态       | 纵向手机屏 9:19.5，1080×2400，Material Design 3 骨架但风格化定制                              | 安卓原生              |
| 底色       | `#0d0805`（近黑棕）/ `#1a1108` / `#1f1610`                                               | 酒馆深夜底            |
| 卡片面     | `linear-gradient(160deg, rgba(34,22,14,0.85), rgba(20,12,8,0.9))`                            | 老旧木板              |
| 描边       | `#3a2616`（暗棕）/ `#5a3a1e`（中棕）                                                       | 木边/铜框             |
| 主品牌色   | `#cb9767` 黄铜色                                                                             | 标题、数值、品牌      |
| 次要文字   | `#d6c0a9` 羊皮纸 / `#8a6a4a` 暗铜 / `#6a5238` 极暗                                       | 层级文字              |
| CTA / 危险 | `linear-gradient(180deg,#9a3827,#7a2a1d)` 暗血红，描边 `#b04a36`                           | 主按钮、游戏中状态    |
| 高亮警示   | `#e94560` 红                                                                                 | 质疑、淘汰、错误      |
| 成功       | `#4ade80`                                                                                    | 出牌、胜利            |
| 章节黄     | `#f5a623`                                                                                    | 小节标题              |
| 圆角       | 8px（按钮/输入）/ 12px（卡片）/ 14px（模态/大卡）/ 20px（pill 标签）                           |                       |
| 阴影       | `0 20px 60px rgba(0,0,0,0.6)` + `inset 0 1px 0 rgba(214,192,169,0.08)`                     | 深而柔                |
| 顶光       | `radial-gradient(ellipse 50% 60% at 50% 0%, rgba(203,151,103,0.18), transparent 70%)`        | 吊灯打光              |
| 暗角       | `radial-gradient(ellipse 75% 60% at 50% 40%, transparent 45%, rgba(0,0,0,0.55) 100%)`        | 四角压暗              |
| 字体       | `"Segoe UI","PingFang SC","Microsoft YaHei"` + 等宽数字                                      |                       |
| 标题       | 28–38px / weight 800 / letter-spacing 2px /`text-shadow: 0 2px 12px rgba(203,151,103,0.35)` | 黄铜发光              |
| 副标题     | 12–14px / letter-spacing 3–4px /`#8a6a4a`                                                  | 宽字距小字            |
| 装饰元素   | 木纹、铜钉、皮革、扑克牌、骷髅、子弹、酒瓶、吊灯、烟雾                                         |                       |
| 交互态     | hover 上浮 1–2px + 描边变`#9a3827` + 暖色光晕                                               |                       |

---

## 1. 通用风格后缀（Prompt Suffix）

> 把下面这段拼接到任意分屏 prompt 的末尾，可保证整套出图风格统一。

```
--ar 9:19.5 --style raw --v 6 --quality 2
mobile UI screen, dark vintage tavern aesthetic, dim warm overhead lamp light, weathered wood textures, brass accents, dark brown #1a1108 background, amber #cb9767 brand color, blood-red #9a3827 accents, parchment #d6c0a9 text, soft vignette, cinematic moody lighting, high detail, Material Design 3 layout骨架, Android phone frame, 9:19.5 aspect ratio, no text artifacts, clean UI mockup, Figma style
```

---

## 2. 分屏提示词

### 2.1 启动 / 闪屏（Splash）

**画面目标**：品牌氛围第一印象，仅 Logo + 标语 + 加载指示。

```
A splash screen for a mobile game called "Liar's Bar" (骗子酒馆), vertical Android phone layout 9:19.5.
Center: an emblem-style logo combining a vintage playing card (Ace of Spades), a brass beer mug, and a small skull motif, embossed brass #cb9767 with subtle glow.
Below the logo: title "LIAR'S BAR" in bold serif, letter-spaced, amber color with soft warm glow text-shadow.
Below title: small subtitle "骗子酒馆" in wide letter-spacing muted brown #8a6a4a.
Bottom: a thin horizontal loading bar in blood-red #9a3827 with brass ticks.
Background: very dark brown #0d0805, faint weathered wood plank texture, a single overhead tavern lamp casting a warm radial pool of light from the top, heavy vignette darkening all four corners.
No other UI elements, no extra text.
```

---

### 2.2 登录（Login）

**画面目标**：单卡居中，账号密码 + 主 CTA + 跳转注册。

```
A dark tavern-themed login screen for a mobile game, vertical Android 9:19.5.
Top 35%: a moody background illustration of a dim vintage bar counter with bottles, a hanging brass lamp, smoke, all heavily darkened by a #0d0805 overlay so only silhouettes show.
Centered card (90% width, rounded 14px): dark wood gradient panel rgba(34,22,14,0.92)→rgba(20,12,8,0.95), 1px brass border #5a3a1e, soft outer shadow, subtle inset top highlight.
Inside the card top: title "LIAR'S BAR" in amber #cb9767 bold serif with warm glow, subtitle "骗子酒馆 · 登录" in muted brown letter-spaced.
Two form fields (username, password): dark inset background rgba(13,8,5,0.6), 1px border #3a2616, rounded 8px, parchment #d6c0a9 placeholder text, focus state shows blood-red #9a3827 border with soft red glow ring.
Primary CTA button full width: blood-red vertical gradient #9a3827→#7a2a1d, 1px border #b04a36, label "进入酒馆" in parchment #f5e6d3 bold, soft red drop shadow.
Below CTA: a small link line "还没有账号？立即注册" in muted brown with amber link.
Top of screen (over status bar): a thin toast notification pill, blood-red border, "用户名或密码错误" warning.
--ar 9:19.5
```

---

### 2.3 注册（Register）

**画面目标**：与登录同骨架，多字段（用户名、邮箱、密码、确认密码、邀请码）。

```
A dark tavern-themed register screen for a mobile game, vertical Android 9:19.5, same visual language as the login screen.
Centered card (92% width, rounded 14px) on a dark #0d0805 wood-texture background with a soft top lamp glow.
Card: title "加入酒馆" in amber #cb9767 with glow, subtitle "REGISTER" in wide letter-spaced muted brown.
Five vertically stacked input fields: 用户名 / 邮箱 / 密码 / 确认密码 / 邀请码(可选).
Inputs: dark inset rgba(13,8,5,0.6), brass-dark border #3a2616, rounded 8px, parchment text, labels in muted cream above each field.
Inline validation hint under the password field in tiny red #e94560: "密码至少 8 位".
Primary CTA: blood-red gradient #9a3827→#7a2a1d, label "创建账号", full width.
Secondary text button below: "已有账号？返回登录" amber link.
A small brass beer-mug icon at the very top of the card as a brand mark.
--ar 9:19.5
```

---

### 2.4 大厅 / 首页（Lobby / Home）

**画面目标**：底部 Tab + 顶部品牌栏 + 三个数据卡 + 三个主操作大按钮 + 活跃房间列表。

```
A home lobby screen for a mobile card game "Liar's Bar", vertical Android 9:19.5, dark tavern aesthetic.
Top app bar: left side brand "Liar's Bar" amber bold + "骗子酒馆" subtitle muted brown; right side a small online-count pill (amber #cb9767 text, dark brown bg, "128 在线"), and a small round avatar circle.
Stats row: three equal cards side by side, each a dark wood gradient panel with 1px brass-dark border, rounded 12px; large amber number on top (活跃房间 12 / 在线玩家 128 / 匹配中 8), small muted label below.
Primary actions: a 2-column grid of large action buttons (rounded 12px, dark wood panel, brass-dark border):
  - "快速匹配" with a crossed-swords icon, hover state shown with blood-red border
  - "创建房间" with a tavern-door icon
  - "个人中心" with a person-silhouette icon
Each button: icon at top in amber, label below in parchment.
Below: section title "活跃房间" in amber letter-spaced, with a small "刷新" icon button on the right.
Room list: vertical stack of room cards (rounded 12px, dark wood gradient, 1px brass-dark border, hover state on the first card with blood-red border + warm overlay). Each card shows: room name (parchment bold), "3 人已准备" small muted text on the left; on the right a player count "3/4", a status pill (waiting = amber outline on translucent amber; playing = blood-red outline on dark red), and a small "加入" button in blood-red gradient.
Bottom navigation bar (Material 3): four tabs with brass icons and parchment labels — 大厅 / 战绩 / 规则 / 我的. Active tab "大厅" highlighted in amber.
Background: full-bleed dark tavern image at 15% opacity behind everything, with #0d0805 base.
--ar 9:19.5
```

---

### 2.5 匹配等待（Match Wait）

**画面目标**：转圈动画 + 实时人数 + 取消按钮 + AI 补位倒计时提示。

```
A matchmaking waiting screen for a mobile card game, vertical Android 9:19.5, dark tavern aesthetic.
Center: a large circular progress ring (brass #cb9767 with a blood-red glow trail), inside the ring an animated count-up number "12s" in amber bold, below it small muted text "正在为您寻找对手...".
Above the ring: a row of 4 player slot avatars — 2 filled with dim player avatars (parchment circles with initials), 2 empty (dashed brass-dark outline with a faint "+" icon). Filled slots have a small amber dot indicator.
Below the ring: a horizontal progress bar showing "Match Progress", brass track with blood-red fill at ~45%.
A small info banner below: "10 秒后将自动加入 AI 玩家补位" in muted amber text with a small skull icon.
Bottom: a single outlined "取消匹配" button in muted brown, with a small back-arrow icon.
Background: dark #0d0805 with very faint silhouettes of a tavern interior (bottles, bar counter), heavy radial vignette focusing attention on the center ring.
Top: small back chevron and title "匹配中".
--ar 9:19.5
```

---

### 2.6 游戏房间 / 牌桌（Game Room）

**画面目标**：核心对战界面。顶部状态条 + 四位玩家头像围绕牌桌 + 中央公共牌堆 + 底部手牌 + 操作按钮 + 侧滑聊天。

```
A game room screen for a multiplayer card bluffing game "Liar's Bar", vertical Android 9:19.5, immersive dark tavern aesthetic.
Top status bar (translucent dark): "← 返回" chevron, center badges — "第 1 轮" / "回合 3" / "目标牌: K" all in small pill shapes (amber text on dark translucent bg), right side "3 人存活" with a tiny skull icon and a "规则" icon button.
Center play area (60% of screen height): a wooden round poker table viewed from above-front, dark green felt center with worn edges, surrounded by a brass-rimmed dark wood rail. On the felt: a small face-down card pile (公共牌堆) in the middle with a "目标牌 K" tag in amber above it.
Four player avatars positioned around the table (top, left, right, and bottom-front = self):
  - Each opponent avatar: a circular portrait with brass border, name below in parchment, three small heart-like HP dots below the name (filled = blood-red #e94560, empty = dark slot), and a fan of face-down card backs showing hand count.
  - The currently active player has a glowing amber ring around their avatar and a small "出牌中…" tag.
  - One opponent has an "AI" tag in tiny red text next to their name.
Bottom player area (self): a horizontal fan of 6 playing cards (face up, showing A/K/Q/J of spades and hearts), each card parchment white with a dark border, slight tilt, the currently selected cards lifted up with a blood-red glow ring.
Bottom action bar (above the bottom nav): three pill buttons in a row:
  - "出牌" success green #4ade80 with dark text
  - "质疑上家" blood-red #e94560 outline on translucent red
  - "过" muted brown outlined
A small chat bubble icon floating at the right edge showing "3" unread badge in red.
Effect overlay shown: a "质疑成功!" big amber title with radial red burst behind it, semi-transparent, mid-animation.
Background: heavy dark tavern interior, hanging lamp above the table casting a warm pool of light, smoke, vignette.
--ar 9:19.5
```

---

### 2.7 个人中心（Profile）

**画面目标**：头像 + 昵称 + ELO + 战绩统计 + 成就/头像选择入口。

```
A profile / personal center screen for a mobile card game, vertical Android 9:19.5, dark tavern aesthetic.
Top: a banner card (rounded 14px, dark wood gradient, brass-dark border) — on the left a large circular avatar (a stylized vintage portrait illustration in sepia), on the right the nickname "Player_001" in parchment bold and below it the ELO score "1248" in large amber #cb9767 with a small ▲ trend arrow in green.
Below the banner: an "编辑资料" small outlined button on the right.
Stats grid (2×3): six stat cards, each a small dark wood panel with brass-dark border, rounded 12px:
  - 总局数 156
  - 胜场 89
  - 败场 67
  - 存活率 57%
  - 撒谎成功率 63%
  - 质疑成功率 48%
Each card: large amber number on top, small muted label below.
Section: "近期成就" in amber letter-spaced. Below: a horizontal scroll row of badge icons (circular brass medallions with engraved symbols — skull, ace, beer mug, bullet), some greyed out / locked.
Section: "头像列表" — a 4-column grid of selectable avatar thumbnails, the currently selected one has an amber ring.
Bottom navigation bar with "我的" tab active in amber.
Background: dark #0d0805, very faint tavern texture, soft top lamp glow.
--ar 9:19.5
```

---

### 2.8 战绩历史（History）

**画面目标**：按时间倒序的战绩列表 + 顶部筛选 Tab（全部 / 胜利 / 失败）。

```
A match history screen for a mobile card game, vertical Android 9:19.5, dark tavern aesthetic.
Top app bar: title "战绩" in amber bold, a small filter icon button on the right.
Below app bar: a horizontal scrollable tab row with three pills — "全部" (active, amber fill), "胜利", "失败" (inactive, muted brown outline).
A summary card at the top: small dark wood panel showing "近 20 场胜率 55%" with a tiny inline bar chart in amber/red.
List of match cards (vertical stack): each card is a dark wood gradient panel, rounded 12px, brass-dark border, with:
  - Left: a vertical color stripe (green for win, blood-red for loss).
  - Center top: "胜利 · 第 1 名" or "失败 · 第 4 名" in parchment bold, with a small placement number badge.
  - Center body: two rows of small muted text — "2026-07-04 21:32" / "12 回合 · 撒谎 4 次 · 质疑 2 次".
  - Right: a small ELO delta chip (+20 green or -15 red) and a chevron icon.
Cards alternate subtly with very slight background variation for readability.
One card mid-list shown in expanded state revealing a small "查看回放" outlined button.
Bottom navigation bar with "战绩" tab active in amber.
Background: dark #0d0805 with faint wood texture.
--ar 9:19.5
```

---

### 2.9 游戏结算（Settlement）

**画面目标**：单局结束后的弹层，展示排名、本局数据、ELO 变化、再战/返回按钮。

```
A match settlement modal overlay for a card bluffing game, vertical Android 9:19.5, dark tavern aesthetic.
Full-screen dim backdrop rgba(0,0,0,0.75) blurring the game room behind.
Centered modal card (88% width, rounded 14px): dark wood gradient panel with 1px brass border #5a3a1e and a soft outer glow, inner top highlight.
Modal top: a large banner row — for the winner a glowing amber crown icon above "胜利！" in amber bold serif; for the loser a skull icon above "淘汰" in blood-red.
Player ranking list (4 rows): each row shows — placement number (1–4) in amber on the left, avatar circle, nickname in parchment, and on the right a small ELO delta chip (green +20 or red -15).
The current player's own row is highlighted with a thin amber border.
Below ranking: a 2-column mini stats grid for the local player — 本局撒谎 3 次 / 质疑 2 次 / 存活 8 回合 / 用时 6:24.
Bottom of modal: two buttons side by side — primary "再来一局" blood-red gradient CTA, secondary "返回大厅" muted brown outlined.
A small "分享战绩" share icon in the top-right corner of the modal.
--ar 9:19.5
```

---

### 2.10 规则说明（Rules）

**画面目标**：底部抽屉或全屏弹层，分章节展示规则，可滚动。

```
A game rules modal screen for a card game "骗子酒馆", vertical Android 9:19.5, dark tavern aesthetic.
Full-screen overlay with rgba(0,0,0,0.75) backdrop.
Centered scrollable card (90% width, 85% height, rounded 14px): dark wood gradient panel, brass-dark border, soft outer shadow.
Card top: a centered title "📖 骗子酒馆 规则" in amber #cb9767 bold, a small close "×" button in the top-right corner in muted brown (hover state shown in blood-red).
Body content as a vertical scroll of sections, each separated by a thin brass-dark divider:
  Section 1 — "目标": small amber #f5a623 heading, parchment body text describing 4-player bluffing game.
  Section 2 — "牌组": heading + parchment text, with a small inline row of 4 card icons (A/K/Q/J) face up.
  Section 3 — "目标牌": heading + text with a highlighted amber "K" chip showing the current target card example.
  Section 4 — "出牌": heading + text, a small inline illustration of 2 face-down cards being played.
  Section 5 — "质疑": heading + text with two example lines — "撒谎 → 受罚" in blood-red, "说真话 → 质疑者受罚" in amber.
  Section 6 — "俄罗斯轮盘惩罚": heading + a small horizontal bullet-row showing 6 bullet chambers, 1 filled red.
  Section 7 — "操作": heading + three small pill buttons "出牌" green / "质疑上家" red / "过" muted.
Subtle brass scroll indicator on the right edge.
--ar 9:19.5
```

---

### 2.11 设置（Settings）

**画面目标**：分组列表 — 账号、音效、通知、关于、退出登录。

```
A settings screen for a mobile card game, vertical Android 9:19.5, dark tavern aesthetic.
Top app bar: "← 设置" back chevron + title in amber bold.
Grouped list sections, each group is a dark wood gradient card (rounded 12px, brass-dark border) with a small amber letter-spaced group title above:
  Group "账号":
    - 修改头像 (right chevron)
    - 修改昵称 (right chevron)
    - 修改密码 (right chevron)
  Group "声音与通知":
    - 音效 (Material 3 toggle switch, on state = amber)
    - 背景音乐 (toggle, off state)
    - 匹配成功通知 (toggle, on state)
  Group "关于":
    - 版本号 v1.0.0 (right side shows version, no chevron)
    - 用户协议 (right chevron)
    - 隐私政策 (right chevron)
Each row: parchment #d6c0a9 label on the left, control/chevron on the right, separated by thin brass-dark dividers, 56dp tap height.
Bottom of screen: a single full-width "退出登录" button in blood-red outline on translucent red, with a small leave icon.
Bottom navigation bar with "我的" tab active.
Background: dark #0d0805, very faint wood texture.
--ar 9:19.5
```

---

## 3. 资产 / 图标提示词

### 3.1 应用图标（App Icon）

```
A mobile app icon for a card-bluffing game "Liar's Bar", square 1024×1024 with rounded mask in mind.
Center: an emblem combining a vintage Ace of Spades card, a brass beer mug, and a small skull at the bottom, all embossed in brass #cb9767 on a dark blood-red #9a3827 shield-like crest background.
Dark wood texture behind the crest, soft inner bevel, slight outer amber glow.
No text, highly detailed, app-store quality, vintage tavern gambling mood.
--v 6 --quality 2
```

### 3.2 品牌 Logo（横向）

```
A horizontal logo lockup for "Liar's Bar", dark tavern aesthetic.
Left: an emblem mark — Ace of Spades crossed with a brass beer mug, small skull below, brass emboss.
Right of mark: wordmark "LIAR'S BAR" in bold vintage serif, brass #cb9767 with soft warm glow, letter-spaced.
Below wordmark: subtitle "骗子酒馆" in muted brown, wide letter-spacing.
Transparent background, centered, high detail.
```

### 3.3 卡牌背面纹理

```
A seamless playing card back design for a dark tavern card game, 2:3 aspect ratio.
Dark burgundy #5a1a1a background with faint wood grain, centered brass #cb9767 emblem of a skull inside an ornate baroque frame, surrounded by vintage filigree and small bullet motifs in the corners.
Symmetrical, intricate line engraving style, slight worn paper texture, no text.
```

### 3.4 头像占位（默认头像）

```
A set of 8 default avatar illustrations for a dark tavern card game, each a 1:1 square.
Vintage sepia-toned illustrated portraits of tavern archetypes: the gambler, the barkeep, the smuggler, the drunkard, the witch, the soldier, the noble, the urchin.
All rendered in muted amber/brown palette with subtle brass border frame, hand-drawn ink engraving style, no text.
```

---

## 4. 使用建议

1. **顺序生成**：建议先生成 §3.1 应用图标和 §3.2 品牌 Logo，确立品牌调性；再生成 §2.1 闪屏；最后按 §2.2 → §2.11 逐屏生成。
2. **保持一致性**：每屏 prompt 末尾务必追加 §1 的「通用风格后缀」，避免出图风格漂移。
3. **文字处理**：图像模型对中文渲染差，可让模型只出"无字骨架"，文字在 Figma / Android Studio 中后续叠字。若用 Midjourney v6 / nano-banana 等强文字模型，可直接出带字版本。
4. **配色复用**：开发阶段直接引用 §0 的色值表入 `colors.xml` / Compose Theme，让最终产物与设计稿一致。
5. **Android 适配**：所有屏幕按 1080×2400 / dp 单位设计；底部导航 80dp 高，状态栏留出 24dp 安全区；触控目标 ≥ 48dp。
6. **暗色主题**：本项目天然是 Dark-only 主题，无需提供 Light 主题变体。

---

## 5. 配色速查（Android `colors.xml` 片段）

```xml
<resources>
  <!-- Base -->
  <color name="bg_base">#FF0D0805</color>
  <color name="bg_panel_top">#FF22160E</color>
  <color name="bg_panel_bottom">#FF140C08</color>
  <color name="border_dark">#FF3A2616</color>
  <color name="border_brass">#FF5A3A1E</color>

  <!-- Brand / Text -->
  <color name="brand_brass">#FFCB9767</color>
  <color name="text_parchment">#FFD6C0A9</color>
  <color name="text_muted">#FF8A6A4A</color>
  <color name="text_dim">#FF6A5238</color>

  <!-- CTA / Status -->
  <color name="cta_red_top">#FF9A3827</color>
  <color name="cta_red_bottom">#FF7A2A1D</color>
  <color name="cta_red_border">#FFB04A36</color>
  <color name="accent_alert">#FFE94560</color>
  <color name="success">#FF4ADE80</color>
  <color name="section_yellow">#FFF5A623</color>
</resources>
```
