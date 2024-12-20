<script setup>
import {reactive} from 'vue'
import * as App from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
})

function greet() {
  App.Greet(data.name).then(result => {
    data.resultText = result
  })
}


</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="greet">Greet</button>
    </div>
  </main>

</template>

<style scoped>
.result {
  height: 24px;
  line-height: 24px;
  margin: 2.5rem auto;
  font-size: 1.2rem;
  color: #2c3e50;
  font-weight: 600;
  text-align: center;
  letter-spacing: 0.5px;
}

.input-box {
  display: flex;
  gap: 20px;
  justify-content: center;
  align-items: center;
  margin: 0 auto;
  max-width: 500px;
  padding: 0 20px;
}

.input-box .btn {
  width: auto;
  min-width: 90px;
  height: 40px;
  line-height: 40px;
  border-radius: 8px;
  border: none;
  padding: 0 20px;
  cursor: pointer;
  background-color: #4CAF50;
  color: white;
  font-weight: 600;
  font-size: 15px;
  transition: all 0.2s ease;
  letter-spacing: 0.5px;
}

.input-box .btn:hover {
  background-color: #45a049;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.2);
}

.input-box .btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(76, 175, 80, 0.1);
}

.input-box .input {
  width: 280px;
  height: 40px;
  line-height: 40px;
  padding: 0 16px;
  border: 2px solid #e8e8e8;
  border-radius: 8px;
  background-color: #fff;
  transition: all 0.3s ease;
  font-size: 15px;
  color: #2c3e50;
}

.input-box .input:hover {
  border-color: #c0c4cc;
  background-color: #fafafa;
}

.input-box .input:focus {
  outline: none;
  border-color: #4CAF50;
  box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.1);
  background-color: #fff;
}
</style>
