<template>
  <Card class="surface-card shadow-4 border-round sm:col-12 md:col-6">
    <template #title>
      <h4 class="m-2">
        {{ msg }}
      </h4>
    </template>
    <template #content>
      <div class="justify-content-center">
        <div class="field ">
          <span class="p-input-icon-left w-full">
            <i class="pi pi-user" />
            <InputText
              id="go-username"
              v-model="username"
              type="username"
              placeholder="Utilisateur"
              class="w-full"
              aria-label="veuillez saisir votre utilisateur"
              aria-describedby="username-help"
            />
            <!-- <small id="username-help" class="p-error">cet utilisateur n'existe pas.</small>-->
          </span>
        </div>
        <div class="field w-full">
          <span class="p-input-icon-left w-full">
            <i class="pi pi-power-off" />
            <Password
              id="go-password"
              v-model="password"
              type="password"
              placeholder="Mot de passe"
              class="w-full"
              input-class="w-full"
              aria-label="veuillez saisir votre mot de passe"
              toggle-mask
              :feedback="false"
            />
          </span>
        </div>
      </div>
    </template>
    <template #footer>
      <div class="justify-content-end text-right">
        <Button label="CONNEXION" class="p-button-raised p-button-info " @click.prevent="" />
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
// import Message from 'primevue/message';
import { getLog } from '../config';

const moduleName = 'LoginUser';
const username = ref('go-admin');
const password = ref('');

/* const urlUserLogin = `${DEFAULT_BASE_SERVER_URL}/data/layers_stats.json`;
const loadedData = ref(false);
const errorDataFetch = ref(null);
const password_hash = ref('');
const sha256Visible = ref(false);
const validLoginForm = ref(false);
const feedbackVisible = ref(false);
const feedbackText = ref('Veuillez vous authentifier SVP.');
const feedbackType = ref('info'); */

const log = getLog(moduleName, 4, 2);

onMounted(() => {
  const method = 'mounted()';
  log.t(`${method}`);
});

defineProps({
  msg: {
    type: String,
    required: true,
  },
});
</script>

<style>
</style>
