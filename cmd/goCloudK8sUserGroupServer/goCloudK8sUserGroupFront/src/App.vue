<template>
  <header>
    <Toolbar class="w-full m-0 p-0 bg-primary-500 border-0">
      <template #start>
        <span class="pl-2 text-white">
          {{ APP_TITLE }} - version : {{ VERSION }}
        </span>
      </template>
      <template #end>
        <Button icon="pi pi-power-off" class="p-button-rounded p-button-sm" title="Logout" />
      </template>
    </Toolbar>
  </header>
  <main>
    <div class="flex">
      <div class="col-12">
        <LoginUser :msg="`Authentification ${APP_TITLE}:`" />
      </div>
    </div>
  </main>
</template>

<script setup>
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import { onMounted } from 'vue';
import LoginUser from './components/LoginUser.vue';
import {
  APP, APP_TITLE, BUILD_DATE, VERSION, getLog,
} from './config';

const log = getLog(APP, 4, 2);

onMounted(() => {
  log.t(' mounted()');
  log.w(`${APP} - ${VERSION}, du ${BUILD_DATE}`);
});

</script>

<style>
html, body {
  padding: 0;
  margin: 0;
}
</style>
