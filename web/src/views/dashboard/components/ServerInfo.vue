<template>
  <Card title="系统信息" class="md:w-2/4 w-full !md:mt-0">
    <div class="py-4 px-4">
      <a-descriptions :column="2">
        <a-descriptions-item label="运行用户"> {{ info.val?.user }} </a-descriptions-item>
        <a-descriptions-item label="运行主机">
          {{ info.val?.hostname }}
        </a-descriptions-item>
        <a-descriptions-item label="系统架构">
          <span v-if="isReady"> {{ info.val.os.goos }}-{{ info.val.os.goarch }} </span>
        </a-descriptions-item>
        <a-descriptions-item label="CPU核心">
          <span v-if="isReady">
            {{ info.val.os.cpu_num }}
          </span>
        </a-descriptions-item>
        <a-descriptions-item label="编译环境">
          <span v-if="isReady">
            {{ info.val.os.compiler }}
          </span>
        </a-descriptions-item>
        <a-descriptions-item label="GO版本">
          <span v-if="isReady">
            {{ info.val.os.go_version }}
          </span>
        </a-descriptions-item>
        <a-descriptions-item label="当前协程数量">
          <span v-if="isReady">
            <CountTo :startVal="1" :endVal="info.val.os.goroutine_num" />
          </span>
        </a-descriptions-item>
        <a-descriptions-item label="CPU使用率">
          <span v-if="isReady"> {{ info.val.cpu.used_percent.toFixed(2) }}% </span>
        </a-descriptions-item>
        <a-descriptions-item label="内存使用率" :span="2">
          <span v-if="isReady">
            {{ info.val.ram.used }}/ {{ info.val.ram.total }} &nbsp; &nbsp;{{
              info.val.ram.used_percent.toFixed(2)
            }}%
          </span>
        </a-descriptions-item>
        <a-descriptions-item label="硬盘使用率" :span="2">
          <span v-if="isReady">
            {{ info.val.disk.used }}/ {{ info.val.disk.total }} &nbsp; &nbsp;
            {{ info.val.disk.used_percent.toFixed(2) }}%
          </span>
        </a-descriptions-item>
      </a-descriptions>
    </div>
  </Card>
</template>
<script lang="ts" setup>
  import { CountTo } from '/@/components/CountTo/index';
  import {
    Card,
    Descriptions as ADescriptions,
    DescriptionsItem as ADescriptionsItem,
    // List as AList,
    // ListItem as AListItem,
    // Row as ARow,
    // Col as ACol,
  } from 'ant-design-vue';
  import { getServerInfo } from '/@/api/common';
  import { ServerInfo } from '/@/api/common/model';
  import { onBeforeUnmount, onMounted, reactive, ref } from 'vue';

  type val = {
    val: ServerInfo;
  };

  const info = reactive<val>({ val: {} } as val);
  const isReady = ref(false);
  const refreshT = ref<NodeJS.Timeout | null>(null);
  const refreshInfoFn = async () => {
    await getServerInfo().then((res) => {
      info.val = res;
      isReady.value = true;
    });
  };
  const getInfo = () => {
    refreshT.value = setTimeout(() => {
      refreshInfoFn().then(getInfo);
    }, 3000);
  };

  onMounted(() => {
    refreshInfoFn().then(getInfo);
  });

  onBeforeUnmount(() => {
    if (refreshT.value) {
      clearTimeout(refreshT.value);
      refreshT.value = null;
    }
  });
</script>
