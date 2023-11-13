<template>
  <div
    v-if="!collapse"
    style="position: absolute; bottom: 10px; padding-left: 16px; color: #aaa; font-size: 12px"
  >
    <a @click="openWindow(SITE_URL)" style="color: #aaa">{{ version.version || '非发行版本' }}</a>
    &nbsp; |&nbsp;
    {{ version.timestamp ? formatToDate(dateUtil(version.timestamp)) : '' }}
  </div>
</template>

<script lang="ts">
  import { defineComponent, reactive } from 'vue';
  import { getVersion } from '/@/api/common';
  import { Version } from '/@/api/common/model';
  import { formatToDate, dateUtil } from '/@/utils/dateUtil';
  import { propTypes } from '/@/utils/propTypes';
  import { SITE_URL } from '/@/settings/siteSetting';
  import { openWindow } from '/@/utils';

  export default defineComponent({
    name: 'Version',
    props: {
      collapse: propTypes.bool,
    },
    setup() {
      const version = reactive<Version>({} as Version);
      getVersion().then((res) => {
        for (let k in res) {
          version[k] = res[k];
        }
      });
      return {
        version,
        dateUtil,
        formatToDate,
        SITE_URL,
        openWindow,
      };
    },
  });
</script>
