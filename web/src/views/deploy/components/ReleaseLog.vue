<template>
  <div>
    <a-card title="" v-loading="loadingRef">
      <Steps
        :current="StepCurrent"
        :serverList="StepServerList"
        :status="StepStatus"
        v-if="StepServerList.length > 0"
      />
      <a-tabs :style="{ minHeight: '400px' }" v-model:activeKey="activeKey">
        <a-tab-pane
          class="text-white bg-gray-700"
          v-for="ss in servers"
          :key="ss.id"
          :tab="ss.host"
        >
          <pre class="mx-1">{{ ss.output }}</pre>
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref, watchEffect, PropType } from 'vue';
  import { PageWrapper } from '/@/components/Page';
  import { Alert, Card as ACard, TabPane as ATabPane, Tabs as ATabs } from 'ant-design-vue';
  import { ListItem as DeployListItem, ReleaseOutput } from '/@/api/deploy/model';
  import { useRoute } from 'vue-router';
  import { DeployStatus, DeployStatusShowMsg } from '/@/enums/fieldEnum';
  import { useWebSocket } from '@vueuse/core';
  import { getDeployConsoleWs, detailDeploy, startDeploy } from '/@/api/deploy';
  import { useUserStore } from '/@/store/modules/user';
  import Steps from './components/Steps.vue';

  const props = defineProps({
    current: { type: Number },
    status: { type: String },
    deployDetail: {
      type: Object as PropType<DeployListItem>,
      default: () => {},
    },
  });

  const userStore = useUserStore();
  const route = useRoute();

  const StepServerList = ref([]);
  const StepCurrent = ref(0);
  const StepStatus = ref('wait');

  const { status, open: wsOpen } = useWebSocket(getDeployConsoleWs(deployId), {
    autoReconnect: false,
    heartbeat: false,
    immediate: false,
    protocols: [userStore.getToken, userStore.getCurrentSpaceId.toString()],
    onError: (_, event) => {
      // alertMsg.type = 'error';
      // alertMsg.msg = event.toString();
    },
    onDisconnected: () => {
      console.log('断开ws连接');
    },
    onMessage: (_, event) => {
      console.log(event.data);
      let data: ReleaseOutput;
      try {
        data = JSON.parse(event.data);
      } catch {
        return;
      }
      servers.value[data.server_id].output += data.data;
    },
  });

  onMounted(() => {
    init(deployId);
  });
</script>
