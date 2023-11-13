<template>
  <BasicModal
    v-bind="$attrs"
    @register="register"
    :title="'项目检测-' + projectDetail?.name"
    width="800px"
    @visible-change="handleVisibleChange"
    :closable="false"
    :keyboard="false"
    :maskClosable="false"
  >
    <Result
      v-if="isSuccessRef"
      status="success"
      title="检测成功"
      sub-title="现在可以自由的对该项目进行发布上线操作了"
    />
    <div
      v-if="!isSuccessRef"
      class="pt-3px pr-3px"
      v-loading="loadingRef"
      loading-tip="'项目检测中...'"
    >
      <a-list item-layout="horizontal" :data-source="detectionList">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta
              :description="item.todo + (item.error ? '  (错误提示:' + item.error + ')' : '')"
            >
              <template #title>
                {{ item.title }}
              </template>
              <template #avatar>
                <CloseCircleOutlined style="color: red; font-size: 18px" />
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
    </div>
  </BasicModal>
</template>
<script lang="ts" setup>
  import { BasicModal, useModalInner } from '/@/components/Modal';
  import { DetectionInfoItem } from '/@/api/project/model';
  import { CloseCircleOutlined } from '@ant-design/icons-vue';
  import { useUserStore } from '/@/store/modules/user';
  import { reactive, ref, nextTick } from 'vue';
  import {
    Result,
    List as AList,
    ListItem as AListItem,
    ListItemMeta as AListItemMeta,
  } from 'ant-design-vue';
  import { useWebSocket } from '@vueuse/core';
  import { getDetectionProjectWs } from '/@/api/project';

  const userStore = useUserStore();
  const loadingRef = ref(true);
  const isSuccessRef = ref(false);
  let servers = {};
  let detectionList = reactive<DetectionInfoItem[]>([]);

  const props = defineProps({
    projectDetail: { type: Object },
  });

  const [register] = useModalInner((data) => {
    console.log('useModalInner', data);
  });

  const { close: wsClose, open: wsOpen } = useWebSocket(
    getDetectionProjectWs(props.projectDetail?.id as number),
    {
      autoReconnect: false,
      immediate: false,
      protocols: [userStore.getToken, userStore.getCurrentSpaceId.toString()],
      onError: () => {
        console.log('onError');
        loadingRef.value = false;
      },
      onDisconnected: () => {
        console.log('onDisconnected');
        loadingRef.value = false;
        if (detectionList.values.length == 0) {
          isSuccessRef.value = true;
        }
      },
      onMessage: (ws, event) => {
        let obj: DetectionInfoItem = JSON.parse(event.data);
        if (obj.server_id) {
          let s = servers[obj.server_id];
          if (s) {
            obj.title = obj.title + '[' + s.user + '@' + s.host + ':' + s.port + ']';
          }
        }
        detectionList.push(obj);
      },
    },
  );
  function handleVisibleChange(v) {
    console.log(222, isSuccessRef.value, loadingRef.value);
    if (v && props.projectDetail) {
      nextTick(() => {
        props.projectDetail?.servers.forEach((element) => {
          servers[element.id] = element;
        });
        wsOpen();
      });
    } else {
      detectionList = reactive<DetectionInfoItem[]>([]);
      isSuccessRef.value = false;
      loadingRef.value = true;
      wsClose();
      console.log(111, isSuccessRef.value, loadingRef.value);
    }
  }
</script>
