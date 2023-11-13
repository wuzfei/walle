<template>
  <PageWrapper contentFullHeight title="发布上线">
    <template #extra>
      <a-button> 返回 </a-button>
      <a-button
        type="primary"
        v-if="deployDetail?.status == DeployStatus.Audit"
        @click="startRelease"
      >
        开始部署
      </a-button>
      <a-button
        type="primary"
        v-if="deployDetail?.status == DeployStatus.Release"
        @click="stopRelease"
      >
        取消部署
      </a-button>
    </template>
    <alert v-if="alertMsg.msg" :type="alertMsg.type">
      <template #description>
        {{ alertMsg.msg }}
      </template>
    </alert>
    <a-card title="" v-loading="loadingRef">
      <Steps :step-list="StepServerList.val" v-if="StepServerList.val.length > 0" />
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
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue';
  import { PageWrapper } from '/@/components/Page';
  import { Alert, Card as ACard, TabPane as ATabPane, Tabs as ATabs } from 'ant-design-vue';
  import { ListItem as DeployListItem, ReleaseOutput } from '/@/api/deploy/model';
  import { useRoute } from 'vue-router';
  import { DeployStatus, DeployStatusShowMsg, AlertMsgType } from '/@/enums/fieldEnum';
  import { useWebSocket } from '@vueuse/core';
  import { getDeployConsoleWs, detailDeploy, startDeploy, stopDeploy } from '/@/api/deploy';
  import { useUserStore } from '/@/store/modules/user';
  import Steps from './components/Steps.vue';
  import { ServerSteps } from './components/props';

  type server = {
    id: number;
    host: string;
    output: string;
  };
  type serverItems = { [key: number]: server };

  const userStore = useUserStore();
  const route = useRoute();
  const loadingRef = ref(false);
  const activeKey = ref<number>(0);
  const deployDetail = ref<DeployListItem | null>(null);

  type StepServerListIfc = {
    val: ServerSteps[];
  };
  const StepServerList = ref<StepServerListIfc>({ val: [] });
  const StepCurrent = ref(0);
  const StepStatus = ref('wait');

  const servers = ref<serverItems>({
    0: { id: 0, host: '127.0.0.1', output: '' },
  });
  const deployId = parseInt(route.params?.id as unknown as string);
  const alertMsg = reactive<{ type: AlertMsgType; msg: string }>({
    type: 'info',
    msg: '',
  });

  const { status, open: wsOpen } = useWebSocket(getDeployConsoleWs(deployId), {
    autoReconnect: false,
    heartbeat: false,
    immediate: false,
    protocols: [userStore.getToken, userStore.getCurrentSpaceId.toString()],
    onError: (_, event) => {
      alertMsg.type = 'error';
      alertMsg.msg = event.toString();
    },
    onDisconnected: () => {
      console.log('断开ws连接');
    },
    onMessage: (_, event) => {
      //console.log(event.data);
      let data: ReleaseOutput;
      try {
        data = JSON.parse(event.data);
      } catch {
        return;
      }
      servers.value[data.server_id].output += data.data;
      changeSteps(data.server_id, data.step, data.over);
    },
  });

  function changeSteps(serverId, step, over) {
    console.log(serverId, step, over);
  }

  function showAlertMsg(msg: string, type: AlertMsgType | undefined = 'error') {
    alertMsg.type = type || 'info';
    alertMsg.msg = msg;
  }

  async function init(id: number) {
    loadingRef.value = true;
    await detailDeploy(id, true)
      .then((res) => {
        deployDetail.value = res;
      })
      .catch((e) => {
        showAlertMsg((e as unknown as Error).toString());
      })
      .finally(() => {
        loadingRef.value = false;
      });

    if (deployDetail.value?.servers.length == 0) {
      showAlertMsg('该上线单发布服务器为空，请检查在执行操作！', 'error');
      return;
    }
    deployDetail.value?.servers.forEach((v) => {
      servers.value[v.id] = {
        id: v.id,
        host: v.host,
        output: '',
      };
      StepServerList.value.val.push({
        server_id: v.id,
        name: v.host,
        current: 0,
        status: 'wait',
      });
    });
    console.log('StepServerList.value.val', StepServerList.value.val);
    let sm = DeployStatusShowMsg(deployDetail.value?.status as number);
    showAlertMsg(sm[0], sm[1]);
    //打开websocket
    if (
      deployDetail.value?.status != DeployStatus.Waiting &&
      deployDetail.value?.status != DeployStatus.Audit
    ) {
      wsOpen();
    }
    //wsOpen()
  }

  function startRelease() {
    loadingRef.value = true;
    startDeploy(deployId, true)
      .then(() => {
        loadingRef.value = false;
        if (status.value != 'OPEN') {
          wsOpen();
        }
      })
      .catch((e) => {
        showAlertMsg(e.toString());
        loadingRef.value = false;
      });
  }

  function stopRelease() {
    loadingRef.value = true;
    stopDeploy(deployId, true)
      .then(() => {
        loadingRef.value = false;
      })
      .catch((e) => {
        showAlertMsg(e.toString());
        loadingRef.value = false;
      });
  }

  onMounted(() => {
    init(deployId);
  });
</script>
