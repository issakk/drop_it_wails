<script setup>
import {h, reactive, ref} from 'vue'
import * as App from '../../wailsjs/go/main/App'
import {ElNotification} from "element-plus";

const value = ref()
const success = (count) => {
  ElNotification({
    title: '整理完毕!',
    message: h('i', {style: 'color: teal'}, 'success 整理了' + count + '个文件'),
  })
}
const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
  treeData: [],
  tableData: [],
  percentage: 0,
  currentPath: "",
})

function listFileInfo(path) {
  console.log("listFiles")
  console.log(path)
  App.ListFileInfo(value.value).then(result => {
    data.tableData = result;
    data.currentPath = path;
  })
}

function drop(path) {
  data.percentage = 0
  if (path === "" || path == null) {
    console.log("null")
    return
  }
  data.percentage = 50
  App.Drop(path).then(result => {
    if (result >= 0) {
      data.percentage = 100
      success(result)
    }
  })
}

function openDialog() {
  console.log("openDialog")
  App.OpenFileDialog().then(result => {
    data.treeData = result;
  })
}

</script>

<template>
  <main>
    <div>
      <div style="margin: 10px">
        <el-space>
          <el-tree-select v-model="value" :data="data.treeData" :render-after-expand="false" check-strictly=true
                          lazy:load="load"
                          @change="listFileInfo"/>

          <el-button type="primary" @click="openDialog">Open Dialog</el-button>
          <el-button type="danger" @click="drop(data.currentPath)">DROP IT!</el-button>

        </el-space>
      </div>

      <div style="margin: 10px">
        <el-progress :percentage="data.percentage" :stroke-width="26" :text-inside="true"/>
      </div>

      <div>
        <el-divider content-position="center"/>
        <el-text class="result" type="info">当前选择路径:{{ data.currentPath }}</el-text>
        <el-table :data="data.tableData" border=border>
          <el-table-column label="Name" prop="name"/>
          <el-table-column label="Date" prop="date"/>
          <el-table-column label="Size" prop="size"/>
        </el-table>
      </div>

    </div>
  </main>

</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}


.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
