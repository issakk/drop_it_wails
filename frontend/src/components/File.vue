<script setup>
import {h, reactive, ref} from 'vue'
import * as App from '../../wailsjs/go/main/App'
import {ElNotification} from "element-plus";
import {Folder, Download, FolderOpened} from '@element-plus/icons-vue'

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
      <el-space size="large" alignment="center">
        <el-tree-select 
          v-model="value" 
          :data="data.treeData" 
          :render-after-expand="false"
          check-strictly="true"
          lazy:load="load"
          @change="listFileInfo"
          placeholder="请选择文件夹"
          style="width: 300px"
        />
        <el-button type="primary" :icon="Folder" @click="openDialog">选择文件夹</el-button>
        <el-button type="danger" :icon="Download" @click="drop(data.currentPath)">开始整理</el-button>
      </el-space>
    </div>

    <div class="file-manager__progress" v-show="data.percentage > 0">
      <el-progress 
        :percentage="data.percentage" 
        :stroke-width="20" 
        :text-inside="true"
        status="success"
      />
    </div>

    <div class="file-manager__content">
      <el-card class="path-card" shadow="hover">
        <template #header>
          <div class="path-header">
            <el-icon><FolderOpened /></el-icon>
            <span>当前路径</span>
          </div>
        </template>
        <el-text class="current-path" type="info">
          {{ data.currentPath || '未选择文件夹' }}
        </el-text>
      </el-card>

      <el-table 
        :data="data.tableData" 
        border 
        stripe
        style="margin-top: 20px"
        v-loading="data.percentage > 0 && data.percentage < 100"
        height="400"
        :virtual-scrolling="true"
        :item-size="40"
      >
        <el-table-column label="文件名" prop="name" min-width="200"/>
        <el-table-column label="修改时间" prop="date" width="180"/>
        <el-table-column label="大小" prop="size" width="120"/>
      </el-table>
    </div>
  </main>
</template>

<style scoped>
.file-manager {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.file-manager__controls {
  margin-bottom: 24px;
  display: flex;
  justify-content: center;
}

.file-manager__progress {
  margin: 24px 0;
  transition: all 0.3s ease;
}

.path-card {
  margin-top: 24px;
}

.path-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-path {
  display: block;
  line-height: 1.5;
  word-break: break-all;
}

:deep(.el-table) {
  --el-table-header-bg-color: var(--el-color-primary-light-8);
  border-radius: 8px;
}

:deep(.el-card) {
  border-radius: 8px;
}
</style>
