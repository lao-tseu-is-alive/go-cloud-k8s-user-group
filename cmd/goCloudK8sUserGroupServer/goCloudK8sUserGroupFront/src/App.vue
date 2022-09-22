<template>
  <header>
    <Toolbar class="w-full m-0 p-0 bg-primary-500 border-0">
      <template #start>
        <span class="pl-2 text-white">{{ `${APP_TITLE} v.${VERSION}` }}</span>
      </template>
      <template #end>
        <template v-if="isUserAuthenticated">
          <Button icon="pi pi-sign-out" class="p-button-rounded p-button-sm" title="Logout" @click="logout" />
        </template>
        <Button icon="pi pi-info-circle" class="p-button-rounded p-button-sm" title="A propos..." @click="aboutInfo" />
      </template>
    </Toolbar>
  </header>
  <main>
    <div class="flex">
      <div class="col-12">
        <FeedBack ref="feedback" :msg="feedbackMsg" :msg-type="feedbackType" :visible="feedbackVisible" />
        <template v-if="isUserAuthenticated">
          <h2>Connexion de {{ getUserName() }} [{{ getUserEmail() }}]</h2>
        </template>
        <template v-else>
          <LoginUser
            :msg="`Authentification ${APP_TITLE}:`"
            :backend="BACKEND_URL"
            :disabled="!isNetworkOk"
            @login-ok="loginSuccess"
            @login-error="loginFailure"
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
import FeedBack from './components/FeedBack.vue';
import {
  getUserId, getUserName, getUserEmail, getUserIsAdmin, getTokenStatus, clearSessionStorage,
} from './components/Login';
import {
  APP, APP_TITLE, BACKEND_URL, BUILD_DATE, VERSION, getLog, HOME,
} from './config';
import { isNullOrUndefined } from './tools/utils';

const log = getLog(APP, 4, 2);
const isUserAuthenticated = ref(false);
const isNetworkOk = ref(true);
const feedback = ref(null);
const feedbackMsg = ref(`${APP_TITLE}, v.${VERSION}`);
const feedbackType = ref('info');
const feedbackVisible = ref(false);
const displayFeedBack = (text, type) => {
  log.t(`displayFeedBack() text:'${text}' type:'${type}'`);
  feedbackType.value = type;
  feedbackMsg.value = text;
  feedbackVisible.value = true;
  feedback.value.displayFeedBack(feedbackMsg, feedbackType);
};
const aboutInfo = () => {
  const appInfo = `${APP_TITLE}, v.${VERSION} ${BUILD_DATE}`;
  if (isUserAuthenticated.value) {
    const userInfo = `${getUserName()} id[${getUserId()}] Admin:${getUserIsAdmin()}`;
    displayFeedBack(`${appInfo} ⇒ 😊 vous êtes "authentifié" comme ${userInfo}`, 'info');
  } else {
    displayFeedBack(`${appInfo} ⇒ vous n'êtes pas encore "authentifié"`, 'info');
  }
  feedbackVisible.value = true;
};

const loginSuccess = (v) => {
  log.t(' loginSuccess()', v);
  isUserAuthenticated.value = true;
  getTokenStatus()
    .then((val) => {
      if (val instanceof Error) {
        log.e('# getTokenStatus() ERROR err: ', val);
        if (val.message === 'Network Error') {
          // displayFeedBack(`Il semble qu'il y a un problème de réseau !${val}`, 'error');
        }
        log.e('# getTokenStatus() ERROR err.response: ', val.response);
        log.w('# getTokenStatus() ERROR err.response.data: ', val.response.data);
        if (!isNullOrUndefined(val.response)) {
          let reason = val.response.data;
          if (!isNullOrUndefined(val.response.data.message)) {
            reason = val.response.data.message;
          }
          log.w(`# getTokenStatus() SERVER SAYS REASON : ${reason}`);
        }
      } else {
        log.l('# getTokenStatus() SUCCESS res: ', val);
      }
    })
    .catch((err) => {
      log.e('# getJwtToken() in catch ERROR err: ', err);
      // displayFeedBack(`Il semble qu'il y a un problème de réseau !${err}`, 'error');
    });
};

const loginFailure = (v) => {
  log.w('loginFailure()', v);
  isUserAuthenticated.value = false;
};

const logout = () => {
  log.t('# IN logout()');
  clearSessionStorage();
  isUserAuthenticated.value = false;
  displayFeedBack('Vous vous êtes déconnecté de l\'application avec succès !', 'success');
  setTimeout(() => {
    window.location.href = HOME;
  }, 2000); // after 2 sec redirect to home page just in case
};

onMounted(() => {
  log.t('mounted()');
  log.w(`${APP} - ${VERSION}, du ${BUILD_DATE}`);

  window.addEventListener('online', () => {
    log.w('ONLINE AGAIN :)');
    isNetworkOk.value = true;
    displayFeedBack('⚡⚡🚀  CONNEXION RESEAU RETABLIE :  😊 vous êtes "ONLINE"  ', 'success');
  });
  window.addEventListener('offline', () => {
    log.w('OFFLINE :((');
    isNetworkOk.value = false;
    displayFeedBack('⚡⚡⚠ PAS DE RESEAU ! ☹ vous êtes "OFFLINE" ', 'error');
  });
});

</script>

<style>
html, body {
  padding: 0;
  margin: 0;
}
</style>
