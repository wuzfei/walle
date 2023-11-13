<template>
  <PageWrapper :class="prefixCls" v-loading="loadingRef" title="项目详情" contentBackground>
    <template #extra>
      <a-button> 返回 </a-button>
      <a-button type="primary" @click="openDetection"> 配置检测 </a-button>
    </template>

    <a-card title="基本信息" :bordered="false" class="mt-5">
      <a-descriptions :column="3">
        <a-descriptions-item label="项目名称">
          {{ projectDetail?.name }}
        </a-descriptions-item>
        <a-descriptions-item label="发布环境">
          {{ projectDetail?.environment_id }}
        </a-descriptions-item>
        <a-descriptions-item label="发布模式"> {{ projectDetail?.repo_mode }} </a-descriptions-item>
        <a-descriptions-item label="仓库类型"> {{ projectDetail?.repo_type }} </a-descriptions-item>
        <a-descriptions-item label="仓库地址" :span="2">
          {{ projectDetail?.repo_url }}
        </a-descriptions-item>

        <a-descriptions-item label="是否需要审核">
          {{ projectDetail?.task_audit == 1 ? '是' : '否' }}
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">
          {{ projectDetail?.created_at }}
        </a-descriptions-item>
        <a-descriptions-item label="最后更新时间">
          {{ projectDetail?.updated_at }}
        </a-descriptions-item>
        <a-descriptions-item label="当前状态">
          {{ projectDetail?.status }}
        </a-descriptions-item>
        <a-descriptions-item
          :span="2"
          :label="projectDetail?.is_include == 1 ? '指定发布文件' : '指定排除文件'"
        >
          {{ projectDetail?.excludes }}
        </a-descriptions-item>
        <a-descriptions-item label="简介说明" :span="3">
          {{ projectDetail?.description }}
        </a-descriptions-item>
      </a-descriptions>
    </a-card>

    <a-card title="服务器配置信息">
      <a-row :gutter="16">
        <a-col :span="6">
          <a-descriptions :column="1">
            <a-descriptions-item label="目标服务器目录">
              {{ projectDetail?.target_root }}
            </a-descriptions-item>
            <a-descriptions-item label="目标服务器路径">
              {{ projectDetail?.target_releases }}
            </a-descriptions-item>
            <a-descriptions-item label="保留版本数量">
              {{ projectDetail?.keep_version_num }}
            </a-descriptions-item>
          </a-descriptions>
        </a-col>
        <a-col>
          <div :class="`${prefixCls}__content`">
            <a-list>
              <a-row :gutter="16">
                <template v-for="item in projectDetail?.servers" :key="item.id">
                  <a-col>
                    <a-tooltip :title="item.description">
                      <a-list-item>
                        <a-card :class="`${prefixCls}__card`">
                          <div :class="`${prefixCls}__card-title`">
                            <Icon class="icon" icon="ant-design:hdd-outlined" />
                            {{ item.user }}@{{ item.host }}:{{ item.port }}
                          </div>
                        </a-card>
                      </a-list-item>
                    </a-tooltip>
                  </a-col>
                </template>
              </a-row>
            </a-list>
          </div>
        </a-col>
      </a-row>
    </a-card>

    <a-card title="执行脚本配置">
      <a-descriptions :column="3">
        <a-descriptions-item label="全局环境变量">
          <a-textarea disabled :value="projectDetail?.task_vars" />
        </a-descriptions-item>
        <a-descriptions-item label="编译前操作命令">
          <a-textarea disabled :value="projectDetail?.prev_deploy" />
        </a-descriptions-item>
        <a-descriptions-item label="编译后操作命令">
          <a-textarea disabled :value="projectDetail?.post_deploy" />
        </a-descriptions-item>
        <a-descriptions-item label="发布前操作命令">
          <a-textarea disabled :value="projectDetail?.prev_release" />
        </a-descriptions-item>
        <a-descriptions-item label="发布后操作命令">
          <a-textarea disabled :value="projectDetail?.post_release" />
        </a-descriptions-item>
      </a-descriptions>
    </a-card>
    <ModalDetection v-if="isReady" @register="register" :projectDetail="projectDetail" />
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { onMounted, ref } from 'vue';
  import { PageWrapper } from '/@/components/Page';
  import Icon from '@/components/Icon/Icon.vue';
  import {
    Card as ACard,
    Descriptions as ADescriptions,
    DescriptionsItem as ADescriptionsItem,
    List as AList,
    ListItem as AListItem,
    Row as ARow,
    Col as ACol,
    Tooltip as ATooltip,
    Textarea as ATextarea,
  } from 'ant-design-vue';
  import { Detail } from '/@/api/project/model';
  import { detailProject } from '/@/api/project';
  import { useRoute } from 'vue-router';
  import { useMessage } from '/@/hooks/web/useMessage';
  import { useModal } from '/@/components/Modal';
  import ModalDetection from './ModalDetection.vue';

  const route = useRoute();
  const { createMessage } = useMessage();
  const [register, { openModal: openDetection }] = useModal();

  const projectId = parseInt(route.params?.id as unknown as string);
  const loadingRef = ref(false);
  const isReady = ref(false);
  const projectDetail = ref<Detail | {}>({});
  const prefixCls = 'project-detail';

  async function init(id: number) {
    try {
      loadingRef.value = true;
      projectDetail.value = await detailProject(id, true);
      isReady.value = true;
    } catch (e) {
      let msg = (e as unknown as Error).toString();
      if (msg == '') {
        msg = '获取项目详情失败';
      }
      createMessage.error(msg);
      return;
    } finally {
      loadingRef.value = false;
    }
  }

  onMounted(() => {
    init(projectId);
  });
</script>
<style lang="less" scoped>
  .desc-wrap {
    padding: 16px;
    background-color: @component-background;
  }

  .project-detail {
    &__card {
      width: 100%;
      margin-bottom: -8px;

      .ant-card-body {
        padding: 16px;
      }

      &-title {
        margin-bottom: 5px;
        color: @text-color;
        font-size: 16px;
        font-weight: 500;

        .icon {
          margin-top: -5px;
          margin-right: 10px;
          font-size: 38px !important;
        }
      }

      &-detail {
        padding-top: 10px;
        padding-left: 30px;
        color: @text-color-secondary;
        font-size: 14px;
      }
    }
  }
</style>
