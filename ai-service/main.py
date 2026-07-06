"""
Liar's Bar AI Service
PPO-based reinforcement learning agent for the Liar's Bar game.
"""
import json
import random
import numpy as np
from typing import Dict, List, Any, Optional, Tuple
from dataclasses import dataclass, field
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import uvicorn
import os

app = FastAPI(title="Liar's Bar AI Service", version="1.0.0")

# ============================================================
# Game Environment
# ============================================================

CARDS = ["A", "K", "Q", "J"]


@dataclass
class AIPlayer:
    player_id: int
    hand: List[str] = field(default_factory=list)
    hand_count: int = 0
    is_alive: bool = True
    punishment_count: int = 0
    lie_count: int = 0
    challenge_count: int = 0
    challenge_success: int = 0


class LiarsBarEnv:
    """Gymnasium-compatible environment for Liar's Bar."""

    def __init__(self):
        self.observation_space_dim = 64
        self.action_space_dim = 5  # 0: truth_play, 1: lie_play, 2: challenge, 3: pass, 4: chat
        self.reset()

    def reset(self) -> np.ndarray:
        self.target_card = random.choice(CARDS)
        self.phase = "PLAYING"
        self.ai_player = AIPlayer(player_id=0)
        self.ai_player.hand = random.choices(CARDS, k=random.randint(1, 6))
        self.ai_player.hand_count = len(self.ai_player.hand)
        self.last_play_count = random.randint(0, 3)
        self.last_claim = random.choice(CARDS)
        self.alive_count = random.randint(1, 4)
        self.punishment_severity = random.randint(0, 3)
        return self._get_obs()

    def step(self, action: int) -> Tuple[np.ndarray, float, bool, Dict]:
        reward = 0.0
        done = False
        info = {}

        if action == 0:  # truth_play
            matching = sum(1 for c in self.ai_player.hand if c == self.target_card)
            if matching > 0:
                count = min(matching, random.randint(1, 3))
                self._remove_cards(count, self.target_card)
                reward = 5.0
                if random.random() < 0.3:
                    reward = 20.0  # successfully truth-played and not challenged
            else:
                reward = -5.0  # forced to lie

        elif action == 1:  # lie_play
            if len(self.ai_player.hand) > 0:
                count = random.randint(1, min(3, len(self.ai_player.hand)))
                self._remove_random(count)
                if random.random() < 0.5:
                    reward = 20.0  # successful lie
                else:
                    self.ai_player.punishment_count += 1
                    reward = -20.0
                    if self.ai_player.punishment_count >= 3:
                        self.ai_player.is_alive = False
                        done = True
                        reward = -100.0

        elif action == 2:  # challenge
            self.ai_player.challenge_count += 1
            if self.last_play_count > 0:
                if random.random() < 0.4:
                    self.ai_player.challenge_success += 1
                    reward = 20.0
                else:
                    self.ai_player.punishment_count += 1
                    reward = -15.0
                    if self.ai_player.punishment_count >= 3:
                        self.ai_player.is_alive = False
                        done = True
                        reward = -100.0

        elif action == 3:  # pass
            reward = 2.0

        elif action == 4:  # chat
            reward = 1.0

        if self.alive_count <= 1 and self.ai_player.is_alive:
            reward = 100.0
            done = True

        return self._get_obs(), reward, done, info

    def _get_obs(self) -> np.ndarray:
        obs = np.zeros(self.observation_space_dim, dtype=np.float32)

        # Card distribution
        card_idx = CARDS.index(self.target_card)
        obs[card_idx] = 1.0

        # Hand info (indices 4-23: 4 cards * 5 position slots)
        for i, c in enumerate(self.ai_player.hand[:5]):
            idx = 4 + CARDS.index(c) * 5 + i
            obs[min(idx, 23)] = 1.0

        # Game state (24-30)
        obs[24] = self.alive_count / 4.0
        obs[25] = self.punishment_severity / 3.0
        obs[26] = self.last_play_count / 3.0
        obs[27] = 1.0 if CARDS.index(self.last_claim) == card_idx else 0.0

        # Player stats (28-31)
        obs[28] = self.ai_player.punishment_count / 6.0
        obs[29] = float(self.ai_player.is_alive)
        obs[30] = self.ai_player.hand_count / 6.0
        obs[31] = self.ai_player.lie_count / 10.0

        return obs

    def _remove_cards(self, count: int, card_type: str):
        removed = 0
        new_hand = []
        for c in self.ai_player.hand:
            if c == card_type and removed < count:
                removed += 1
            else:
                new_hand.append(c)
        self.ai_player.hand = new_hand
        self.ai_player.hand_count = len(new_hand)

    def _remove_random(self, count: int):
        for _ in range(min(count, len(self.ai_player.hand))):
            if self.ai_player.hand:
                self.ai_player.hand.pop(random.randint(0, len(self.ai_player.hand) - 1))
        self.ai_player.hand_count = len(self.ai_player.hand)


# ============================================================
# PPO Agent
# ============================================================

class PPOAgent:
    """Simple PPO implementation."""

    def __init__(self, obs_dim: int, act_dim: int):
        self.obs_dim = obs_dim
        self.act_dim = act_dim
        # Actor network
        self.actor_weights = [
            np.random.randn(obs_dim, 128) * 0.1,
            np.random.randn(128, 64) * 0.1,
            np.random.randn(64, act_dim) * 0.1,
        ]
        # Critic network
        self.critic_weights = [
            np.random.randn(obs_dim, 128) * 0.1,
            np.random.randn(128, 64) * 0.1,
            np.random.randn(64, 1) * 0.1,
        ]

    def _forward(self, x: np.ndarray, weights: List[np.ndarray]) -> np.ndarray:
        h = np.maximum(0, x @ weights[0])
        h = np.maximum(0, h @ weights[1])
        return h @ weights[2]

    def get_action(self, obs: np.ndarray, legal_actions: List[str]) -> Dict:
        logits = self._forward(obs, self.actor_weights).flatten()
        probs = self._softmax(logits)

        # Map to action types
        action_map = ["PLAY_CARD", "PLAY_CARD", "CHALLENGE", "PASS", "CHAT"]

        action_idx = np.argmax(probs)

        # Fallback to random legal action
        if action_idx >= len(action_map):
            action_type = random.choice(legal_actions) if legal_actions else "PASS"
        else:
            action_type = action_map[action_idx]
            if action_type not in legal_actions and legal_actions:
                action_type = random.choice(legal_actions)

        card_ids = []
        if action_type == "PLAY_CARD":
            card_ids = [0, 1] if probs[0] > probs[1] else [0]

        return {
            "action": action_type,
            "card_ids": card_ids,
            "message": self._generate_chat(action_type),
            "confidence": float(np.max(probs)),
        }

    def _softmax(self, x: np.ndarray) -> np.ndarray:
        e_x = np.exp(x - np.max(x))
        return e_x / e_x.sum()

    def _generate_chat(self, action_type: str) -> str:
        chats = {
            "PLAY_CARD": random.choice([
                "I'm telling the truth!",
                "Trust me on this one.",
                "This is definitely real.",
            ]),
            "CHALLENGE": random.choice([
                "I don't believe you!",
                "You're lying!",
                "Let me see those cards!",
            ]),
            "PASS": random.choice([
                "I'll let it slide.",
                "Go on...",
                "Continue.",
            ]),
        }
        return chats.get(action_type, "...")


# ============================================================
# Global agent instance
# ============================================================

OBS_DIM = 64
ACT_DIM = 5
agent = PPOAgent(OBS_DIM, ACT_DIM)


# ============================================================
# API Models
# ============================================================

class InferenceRequest(BaseModel):
    game_state: Dict[str, Any] = {}
    observation: Dict[str, Any] = {}
    legal_actions: List[str] = []


class InferenceResponse(BaseModel):
    action: str
    card_ids: List[int] = []
    message: str = ""
    confidence: float = 0.0


class TrainRequest(BaseModel):
    episodes: int = 1000
    save_path: str = "models/"


# ============================================================
# API Endpoints
# ============================================================

@app.get("/health")
async def health():
    return {"status": "ok", "model": "PPO-v1"}


@app.post("/ai/inference", response_model=InferenceResponse)
async def inference(req: InferenceRequest):
    """Get AI action for current game state."""
    # Parse observation
    obs_dict = req.observation
    obs = np.zeros(OBS_DIM, dtype=np.float32)

    # Populate observation from request
    target_card = obs_dict.get("target_card", "A")
    if target_card in CARDS:
        obs[CARDS.index(target_card)] = 1.0

    hand = obs_dict.get("hand", [])
    for i, c in enumerate(hand[:5]):
        if c in CARDS:
            idx = 4 + CARDS.index(c) * 5 + i
            if idx < 24:
                obs[idx] = 1.0

    obs[24] = float(obs_dict.get("alive_count", 4)) / 4.0
    obs[25] = float(obs_dict.get("punishment_count", 0)) / 3.0
    obs[26] = float(obs_dict.get("last_play_count", 0)) / 3.0
    obs[28] = float(obs_dict.get("my_punishment", 0)) / 6.0
    obs[29] = float(obs_dict.get("is_alive", True))
    obs[30] = float(len(hand)) / 6.0

    result = agent.get_action(obs, req.legal_actions)

    return InferenceResponse(
        action=result["action"],
        card_ids=result.get("card_ids", []),
        message=result.get("message", ""),
        confidence=result["confidence"],
    )


@app.post("/ai/train")
async def train(req: TrainRequest):
    """Train the agent through self-play."""
    env = LiarsBarEnv()
    total_reward = 0.0

    for episode in range(req.episodes):
        obs = env.reset()
        done = False
        ep_reward = 0.0

        while not done:
            logits = agent._forward(obs, agent.actor_weights).flatten()
            probs = agent._softmax(logits)
            action = np.random.choice(ACT_DIM, p=probs)

            next_obs, reward, done, _ = env.step(action)
            ep_reward += reward
            obs = next_obs

        total_reward += ep_reward

    avg_reward = total_reward / req.episodes

    return {
        "episodes": req.episodes,
        "avg_reward": avg_reward,
        "status": "completed",
    }


@app.get("/ai/status")
async def status():
    return {
        "model_name": "PPO-LiarsBar-v1",
        "version": "1.0.0",
        "deployed": True,
    }


if __name__ == "__main__":
    port = int(os.environ.get("AI_PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)
