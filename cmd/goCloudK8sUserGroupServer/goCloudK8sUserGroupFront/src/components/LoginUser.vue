<template>
  <Card class="surface-card shadow-4 border-round sm:col-12 md:col-6  md:col-offset-3">
    <template #title>
      <h5 class="m-2">
        {{ msg }}
      </h5>
    </template>
    <template #content>
      <div class="field ">
        <span class="p-input-icon-left w-full">
          <i class="pi pi-user" />
          <InputText
            id="go-username"
            ref="userInput"
            v-model="username"
            type="username"
            placeholder="Utilisateur"
            class="w-full"
            aria-label="veuillez saisir votre utilisateur"
            autofocus
            :disabled="disabled"
            @keyup.enter="onEnterKey"
          />
        </span>
      </div>
      <div class="field w-full">
        <span class="p-input-icon-left w-full">
          <i class="pi pi-power-off" />
          <Password
            id="go-password"
            ref="pwdInput"
            v-model="password"
            placeholder="Mot de passe"
            class="w-full"
            input-class="w-full"
            aria-label="veuillez saisir votre mot de passe"
            toggle-mask
            :feedback="false"
            :disabled="disabled"
            @keyup.enter="onEnterKey"
          />
        </span>
      </div>
      <div v-if="feedbackVisible">
        <Message :severity="feedbackType" class="w-full">
          {{ feedbackText }}
        </Message>
      </div>
    </template>
    <template #footer>
      <div class="justify-content-end text-right">
        <Button
          label="CONNEXION"
          class="p-button-raised p-button-info "
          :disabled="disabled"
          @click.prevent="getJwtToken"
        />
      </div>
    </template>
  </Card>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import Button from 'primevue/button';
import Card from 'primevue/card';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Message from 'primevue/inlinemessage';
import { getLog } from '../config';
import { getPasswordHash, getToken, clearSessionStorage } from './Login';
import { isNullOrUndefined } from '../tools/utils';

const moduleName = 'LoginUser';
const username = ref('');
const password = ref('');
const userInput = ref(null);
const pwdInput = ref(null);
const feedbackVisible = ref(true);
const feedbackText = ref('Veuillez vous authentifier SVP.');
const feedbackType = ref('info');

const log = getLog(moduleName, 4, 2);

const props = defineProps({
  msg: {
    type: String,
    required: true,
    default: 'Authentification',
  },
  backend: {
    type: String,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: null,
  },
});

const emit = defineEmits(['login-ok', 'login-error']);

const displayFeedBack = (text, type) => {
  const validTypes = ['success', 'info', 'warn', 'error'];
  if (validTypes.includes(type)) {
    feedbackType.value = type;
  } else {
    feedbackType.value = 'info';
  }
  feedbackText.value = text;
  feedbackVisible.value = true;
};
const resetFeedBack = () => {
  feedbackText.value = '';
  feedbackType.value = 'info';
  feedbackVisible.value = false;
};
const isValidForm = () => {
  if (username.value.trim().length < 1) {
    displayFeedBack('Veuillez saisir votre utilisateur, il est obligatoire!', 'warn');
    userInput.value.$el.focus();
    return false;
  }
  if (password.value.trim().length < 1) {
    displayFeedBack('Veuillez saisir votre mot de passe, il est obligatoire!', 'warn');
    pwdInput.value.$el.focus();
    return false;
  }
  return true;
};

const getJwtToken = () => {
  log.t('# IN getJwtToken()');
  if (isValidForm()) {
    resetFeedBack();
    try {
      const res = getToken(
        props.backend,
        username.value,
        getPasswordHash(password.value),
      )
        .then((val) => {
          if (val instanceof Error) {
            log.e('# getJwtToken() ERROR err: ', val);
            if (val.message === 'Network Error') {
              displayFeedBack(`Il semble qu'il y a un problème de réseau !${val}`, 'error');
            }
            log.e('# getJwtToken() ERROR err.response: ', val.response);
            log.w('# getJwtToken() ERROR err.response.data: ', val.response.data);
            if (!isNullOrUndefined(val.response)) {
              let reason = val.response.data;
              if (!isNullOrUndefined(val.response.data.message)) {
                reason = val.response.data.message;
              }
              log.w(`# getJwtToken() SERVER SAYS REASON : ${reason}`);
              if ((reason.match(/wrong password/gi) !== null)
                    || (reason.match(/no records found/gi) !== null)) {
                displayFeedBack('Vos informations de connexions sont erronées !', 'warn');
              } else {
                displayFeedBack(`Erreur serveur : ${reason}`, 'error');
              }
            } else {
              displayFeedBack(`ERREUR SERVEUR :  ${val}`, 'error');
            }
            emit('login-error', 'LOGIN FAILED', val);
          } else {
            log.l('# getJwtToken() SUCCESS res: ', val);
            displayFeedBack('Connexion réussie !', 'success');
            emit('login-ok', 'LOGIN SUCCESS', val);
          }
        })
        .catch((err) => {
          log.e('# getJwtToken() in catch ERROR err: ', err);
          displayFeedBack(`Il semble qu'il y a un problème de réseau !${err}`, 'error');
          emit('login-error', 'LOGIN ERROR', err);
        });
      log.l('# getJwtToken() after getToken res:', res);
    } catch (e) {
      log.t('# getJwtToken() TRY CATCH ERROR : ', e);
    }
  } else {
    log.w('le formulaire de connexion est invalide');
  }
  log.t('# GOING OUT getJwtToken()');
};
const onEnterKey = () => {
  log.t('# IN onEnterKey()');
  if (isValidForm()) {
    getJwtToken();
    return true;
  }
  return false;
};

onMounted(() => {
  const method = 'mounted()';
  log.t(`${method}`);
  clearSessionStorage();
});
</script>
