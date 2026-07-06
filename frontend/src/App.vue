<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background: #1a1a2e;
  color: #e0e0e0;
  min-height: 100vh;
}

#app {
  min-height: 100vh;
}

button {
  cursor: pointer;
  border: none;
  border-radius: 8px;
  padding: 10px 24px;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.2s;
}

button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
}

button:active {
  transform: translateY(0);
}

input {
  padding: 12px 16px;
  border: 2px solid #333;
  border-radius: 8px;
  background: #16213e;
  color: #e0e0e0;
  font-size: 16px;
  outline: none;
  transition: border-color 0.2s;
}

input:focus {
  border-color: #e94560;
}

a {
  color: #e94560;
  text-decoration: none;
}
</style>
