<style>
</style>

<template>
  <Card class="surface-card shadow-4 border-round sm:w-12 md:w-6">
    <template #title>
      <h4 class="m-2">{{ msg }}</h4>
    </template>
    <template #content>
      <div class="justify-content-center"></div>
      <div class="field ">
          <span class="p-input-icon-left">
            <i class="pi pi-user"/>
          <InputText id="go-username" type="username" placeholder="Utilisateur"
                     aria-describedby="username-help" class="p-invalid"/>
            <!--          <small id="username-help" class="p-error">cet utilisateur n'existe pas.</small>-->
          </span>
      </div>
      <div class="field">
          <span class="p-input-icon-left">
            <i class="pi pi-power-off"/>
            <Password type="password" id="go-password" placeholder="Mot de passe"
                      toggleMask :feedback="false"></Password>
          </span>
      </div>
    </template>
    <template #footer>
      <Button label="CONNEXION" class="p-button-raised p-button-info right-0"/>
    </template>
  </Card>
</template>

<script setup>
import {onMounted, ref} from "vue";
import Button from "primevue/button";
import Card from "primevue/card";
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import {FilterMatchMode} from "primevue/api";
import {fetchData} from "./fetchData";
import {DEFAULT_BASE_SERVER_URL} from "../config";

const moduleName = "LoginUser";
const urlUserLogin = `${DEFAULT_BASE_SERVER_URL}/data/layers_stats.json`;
const loadedData = ref(false);
const errorDataFetch = ref(null);
const dataUsers = ref([
  {
    username: "aucune données n'est disponible en ce moment...",
  },
]);
const filters = ref({
  global: {value: null, matchMode: FilterMatchMode.CONTAINS},
  layer: {value: null, matchMode: FilterMatchMode.CONTAINS},
});

onMounted(() => {
  const method = "onMounted";
  console.log(`##-->${moduleName}::${method}`);
});

defineProps({
  msg: {
    type: String,
    required: true,
  },
});
</script>
