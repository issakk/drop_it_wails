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
      listFileInfo(path)
    }
  })
}

function openDialog() {
  console.log("openDialog")
  App.OpenFileDialog().then(result => {
    data.treeData = result;
    if (result && result.length > 0) {
      value.value = result[0].value;
      listFileInfo(result[0].value);
    }
  })
}

</script>

<template>
  <main class="file-manager">
    <div class="file-manager__controls">
      <el-space>
        <el-tree-select 
          v-model="value" 
          :data="data.treeData" 
          :render-after-expand="false"
          check-strictly="true"
          lazy:load="load"
          @change="listFileInfo"
        />
        <el-button type="primary" @click="openDialog">Open Dialog</el-button>
        <el-button type="danger" @click="drop(data.currentPath)">DROP IT!</el-button>
      </el-space>
    </div>

    <div class="file-manager__progress">
      <el-progress 
        :percentage="data.percentage" 
        :stroke-width="26" 
        :text-inside="true"
      />
    </div>

    <div class="file-manager__content">
      <el-divider content-position="center"/>
      <el-text class="current-path" type="info">
        当前选择路径:{{ data.currentPath }}
      </el-text>
      <el-table :data="data.tableData" border>
        <el-table-column label="Name" prop="name"/>
        <el-table-column label="Date" prop="date"/>
        <el-table-column label="Size" prop="size"/>
      </el-table>
    </div>
  </main>
</template>

<style scoped>
.file-manager {
  padding: 20px;
}

.file-manager__controls {
  margin-bottom: 20px;
}

.file-manager__progress {
  margin-bottom: 20px;
}

.current-path {
  display: block;
  margin: 1.5rem 0;
  line-height: 20px;
}

/* 删除未使用的 input-box 相关样式 */
</style>
