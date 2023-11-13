<template>
  <Card>
    <Row type="flex" justify="space-around" align="middle">
      <Col :span="10">
        <a-steps :current="currentStep" :status="currentStatus" size="small">
          <template v-for="v in localSteps" :key="v">
            <a-step :title="v.title" :disabled="true" />
          </template>
        </a-steps>
      </Col>

      <Col :span="2">
        <div ref="linesRef">
          <template v-for="(v, i) in stepList" :key="v.server_id">
            <div
              v-if="v.server_id > 0"
              class="line"
              :id="'server-' + v.server_id"
              :style="{
                '--height': lineStyle.val[i]?.height,
                '--top': lineStyle.val[i]?.top,
                '--rotateZ': lineStyle.val[i]?.rotateZ,
                '--heightBefore': lineStyle.val[i]?.heightBefore,
              }"
            >
            </div>
          </template>
          <!-- <div
            class="line"
            v-for="(v, i) in stepList"
            :key="v.server_id"
            :id="'server-' + v.server_id"
            :style="{
              '--height': lineStyle.val[i]?.height,
              '--top': lineStyle.val[i]?.top,
              '--rotateZ': lineStyle.val[i]?.rotateZ,
              '--heightBefore': lineStyle.val[i]?.heightBefore,
            }"
          >
          </div> -->
        </div>
      </Col>
      <Col :span="12">
        <div ref="serverStepsRef">
          <template v-for="ss in stepList" :key="ss.server_id">
            <div class="py-2">
              <a-steps
                :current="ss.current || 0"
                :initial="localSteps.length"
                size="small"
                :status="ss.status"
              >
                <template v-for="(v, i) in serverSteps" :key="v">
                  <a-step
                    :title="v.title + (i == serverSteps.length - 1 ? ':' + ss.name : '')"
                    :disabled="true"
                  />
                </template>
              </a-steps>
            </div>
          </template>
        </div>
      </Col>
    </Row>
  </Card>
</template>

<script lang="ts" setup>
  import { Steps as ASteps, Step as AStep, Row, Col, Card } from 'ant-design-vue';
  import { nextTick, onMounted, ref, PropType, watchEffect } from 'vue';
  import { type Recordable } from '@vben/types';
  import { ServerSteps, statusType } from './props';

  type sty = {
    height: string;
    rotateZ: string;
    top: string;
    heightBefore?: string;
  };
  type vSty = {
    val: { [key: number]: sty };
  };
  const localSteps = [{ title: '前置操作' }, { title: '检出代码' }, { title: '编译打包' }];
  const serverSteps = [{ title: '部署前置操作' }, { title: '部署程序' }, { title: '部署后置操作' }];

  const props = defineProps({
    stepList: {
      type: Array as PropType<(ServerSteps & Recordable<any>)[]>,
      default: () => [],
    },
  });

  const serverStepsRef = ref<HTMLElement | null>(null);
  const linesRef = ref<HTMLElement | null>(null);
  const lineStyle = ref<vSty>({ val: {} });

  const currentStep = ref<number>(0);
  const currentStatus = ref<statusType>('wait');

  watchEffect(() => {
    if (props.stepList) {
      console.log('props.stepList', props.stepList);
      props.stepList.forEach((v) => {
        if (v.server_id == 0) {
          currentStep.value = v.current;
          currentStatus.value = v.status;
        }
      });
    }
  });

  onMounted(() => {
    nextTick(() => {
      console.log('props.stepList', props.stepList);
      let cc = props.stepList.length >= 1 ? props.stepList.length : 0;
      if (cc == 0) {
        return;
      }
      let width = linesRef.value?.offsetWidth || 0;
      let H = serverStepsRef.value?.offsetHeight || 0;
      let halfH = H / 2;
      let oneH = H / cc;
      let lines = linesRef.value?.getElementsByClassName('line');
      let lineHalfNum = Math.floor((lines?.length || 0) / 2);
      let hasOne = (lines?.length || 0) % 2 == 1;
      let scaleRate = 0.8;
      if (lines?.length == 0) {
        return;
      }
      for (let i = 1; i <= lineHalfNum; i++) {
        let _h = halfH - oneH * (i - 0) + oneH / 2;
        //边长
        let borderWidth = Math.floor(Math.sqrt(Math.pow(_h, 2) + Math.pow(width, 2)));
        //缩放比例
        let scale = (borderWidth / width) * scaleRate;
        //旋转的角度
        let dd = (180 * Math.atan(_h / width)) / Math.PI;
        lineStyle.value.val[i - 1] = {
          height: `${_h}px`,
          heightBefore: `${_h / 2}px`,
          top: `-${_h}px`,
          rotateZ: `rotateZ(${-dd}deg) scale(${scale})`,
        };
        lineStyle.value.val[cc - i] = {
          height: `${_h}px`,
          heightBefore: `${_h / 2}px`,
          top: `0px`,
          rotateZ: `rotateZ(${dd}deg) scale(${scale})`,
        };
        if (hasOne && i == lineHalfNum) {
          lineStyle.value.val[lineHalfNum] = {
            height: `0px`,
            heightBefore: `0px`,
            top: '-1px',
            rotateZ: `rotate(0deg) scale(${scaleRate})`,
          };
        }
      }
    });
  });
</script>

<style lang="css" scoped>
  .line::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    box-sizing: border-box;
    width: 100%;
    height: var(--heightBefore);
    transform: var(--rotateZ);
    transform-origin: bottom center;
    border-bottom: 1px solid blue;
  }

  .line {
    position: absolute;
    top: var(--top);
    box-sizing: border-box;
    width: 100%;
    height: var(--height);
  }
</style>
