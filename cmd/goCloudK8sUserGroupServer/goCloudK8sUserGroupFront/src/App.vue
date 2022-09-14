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
        <template v-if="isUserAuthenticated">
          <h2>Connexion de {{ getUserName() }} [{{ getUserEmail() }}]</h2>
        </template>
        <template v-else>
          <LoginUser
            :msg="`Authentification ${APP_TITLE}:`"
            :backend="BACKEND_URL"
            @loginOK="connected"
            @loginError="failedConnect"
          />
        </template>
      </div>
    </div>
  </main>
</template>

<script setup>
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import { onMounted, ref } from 'vue';
import LoginUser from './components/LoginUser.vue';
import { getUserName, getUserEmail } from './components/Login';
import {
  APP, APP_TITLE, BACKEND_URL, BUILD_DATE, VERSION, getLog,
} from './config';

const log = getLog(APP, 4, 2);
const isUserAuthenticated = ref(false);
const connected = (v) => {
  log.t(' connected()', v);
  isUserAuthenticated.value = true;
};

const failedConnect = (v) => {
  log.w('FailedConnect()', v);
};

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
