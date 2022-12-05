<template>
  <div class="w-full card">
    <DataTable
      ref="dt"
      v-model:filters="filters"
      :value="dataUsers"
      removable-sort
      responsive-layout="scroll"
      striped-rows
    >
      <template #header>
        <Toolbar class="m-0 p-0">
          <template #start>
            <span class="p-input-icon-left ">
              <i class="pi pi-search" />
              <InputText v-model="filters['global'].value" placeholder="filtre..." />
            </span>
            <h4 class="m-2 hidden md:inline">
              Liste des {{ getNumUsers }} utilisateurs
            </h4>
          </template>
          <template #end>
            <Button label="Nouvel Utilisateur" icon="pi pi-plus" class="p-button-success mr-0" @click="openNew" />
          </template>
        </Toolbar>
      </template>
      <ColumnGroup type="header">
        <Row>
          <Column field="id" header="id" :sortable="true" />
          <Column field="username" header="Utilisateur" :sortable="true" />
          <Column field="name" header="Nom" :sortable="true" />
          <Column field="is_admin" header="Admin?" :sortable="true" />
          <Column field="is_locked" header="Locked?" :sortable="true" />
          <Column header="Actions" />
        </Row>
      </ColumnGroup>
      <Column field="id" />
      <Column field="username" class="align-right" />
      <Column field="name" class="align-right" />
      <Column field="is_admin" class="align-center" />
      <Column field="is_locked" class="align-center" />
      <Column :exportable="false" style="min-width:8rem">
        <template #body="slotProps">
          <Button icon="pi pi-pencil" class="p-button-rounded  mr-2" @click="editUser(slotProps.data)" />
          <Button icon="pi pi-trash" class="p-button-rounded " @click="confirmDeleteUser(slotProps.data)" />
        </template>
      </Column>
    </DataTable>
  </div>
  <!-- BEGIN DIALOG EDIT USER -->
  <!--
  default width is set to 60vw and below 961px, width would be 75vw
   and finally below 641px width becomes 100%
  -->
  <Dialog
    v-model:visible="userDialog"
    :breakpoints="{ '960px': '75vw', '640px': '100vw' }"
    :style="{ width: '60vw' }"
    :header="`Utilisateur id[${dataCurrentUser.id}], création:${dataCurrentUser.create_time}`"
    :modal="true"
    class="p-fluid"
  >
    <div class="field">
      <label for="name">Nom</label>
      <InputText id="name" v-model.trim="dataCurrentUser.name" required="true" autofocus :class="{ 'p-invalid': submitted && !dataCurrentUser.name }" />
      <small v-if="submitted && !dataCurrentUser.name" class="p-error">Name est obligatoire.</small>
    </div>
    <div class="field">
      <label for="username">Username</label>
      <InputText id="username" v-model.trim="dataCurrentUser.username" required="true" :class="{ 'p-invalid': submitted && !dataCurrentUser.username }" />
      <small v-if="submitted && !dataCurrentUser.username" class="p-error">Username est obligatoire.</small>
    </div>
    <div class="field">
      <label for="password">Mot de passe</label>
      <InputText id="password" v-model.trim="dataCurrentUser.password" type="password" :class="{ 'p-invalid': submitted && isNewUser && !dataCurrentUser.password }" />
      <small v-if="submitted && isNewUser && !dataCurrentUser.password" class="p-error">Le mot de passe est obligatoire.</small>
    </div>
    <div class="field">
      <label for="email">Email</label>
      <InputText id="email" v-model.trim="dataCurrentUser.email" required="true" :class="{ 'p-invalid': submitted && !dataCurrentUser.email }" />
      <small v-if="submitted && !dataCurrentUser.email" class="p-error">L'email est obligatoire.</small>
    </div>
    <div class="field">
      <label for="orgunit">OrgUnit</label>
      <InputText id="orgunit" v-model.trim="dataCurrentUser.orgunit_id" />
    </div>
    <div class="field">
      <label for="groups">Groupes</label>
      <MultiSelect
        v-model="dataCurrentUser.groups_id"
        :options="groupsList"
        option-label="name"
        option-value="id"
        placeholder="Choisissez les groupes"
      />
      <InputText id="groups" v-model.trim="dataCurrentUser.groups_id" />
    </div>
    <div class="field">
      <label for="phone">Téléphone</label>
      <InputText id="phone" v-model.trim="dataCurrentUser.phone" />
    </div>
    <div class="field">
      <label for="description">Commentaire</label>
      <Textarea id="description" v-model="dataCurrentUser.comment" required="true" rows="3" cols="20" />
    </div>

    <div class="formgrid grid">
      <div class="field col">
        <label for="isadmin">Administrateur ?</label>
        <ToggleButton
          v-model="dataCurrentUser.is_admin"
          on-label="Oui"
          off-label="Non"
          on-icon="pi pi-check"
          off-icon="pi pi-times"
          class="w-full sm:w-10rem"
          aria-label="compte administrateur"
        />
      </div>
      <div class="field col">
        <label for="islocked">verrouillé ?</label>
        <ToggleButton
          v-model="dataCurrentUser.is_locked"
          on-label="Oui"
          off-label="Non"
          on-icon="pi pi-check"
          off-icon="pi pi-times"
          class="w-full sm:w-10rem"
          aria-label="compte verrouillé"
        />
      </div>
      <div class="field col">
        <label for="isactive">Compte actif ? </label>
        <ToggleButton
          v-model="dataCurrentUser.is_active"
          on-label="Oui"
          off-label="Non"
          on-icon="pi pi-check"
          off-icon="pi pi-times"
          class="w-full sm:w-10rem"
          aria-label="compte actif"
        />
      </div>
    </div>
    <template #footer>
      <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="hideDialog" />
      <Button label="Save" icon="pi pi-check" class="p-button-text" @click="saveUser" />
    </template>
  </Dialog>
  <!-- END DIALOG EDIT USER -->
  <Dialog v-model:visible="deleteUserDialog" :style="{ width: '450px' }" header="Confirm" :modal="true">
    <div class="confirmation-content">
      <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
      <span v-if="dataCurrentUser">
        Are you sure you want to delete <b>{{ dataCurrentUser.name }}</b>?
      </span>
    </div>
    <template #footer>
      <Button label="No" icon="pi pi-times" class="p-button-text" @click="deleteUserDialog = false" />
      <Button label="Yes" icon="pi pi-check" class="p-button-text" @click="deleteUser" />
    </template>
  </Dialog>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue';
import { useToast } from 'primevue/usetoast';
import Button from 'primevue/button';
import Column from 'primevue/column';
import ColumnGroup from 'primevue/columngroup';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import MultiSelect from 'primevue/multiselect';
import Row from 'primevue/row';
import Toolbar from 'primevue/toolbar';
import Textarea from 'primevue/textarea';
import ToggleButton from 'primevue/togglebutton';
import { FilterMatchMode } from 'primevue/api';
import user from './User';
import { getPasswordHash, getUserId } from './Login';
import { getLog } from '../config';
import { isNullOrUndefined } from '../tools/utils';
import group from './Group';

const moduleName = 'ListUsers';
const timeToDisplayError = 7000;
const timeToDisplaySucces = 4000;

const log = getLog(moduleName, 4, 2);
const loadedData = ref(false);
const loadingData = ref(false);
const loadingGroupsList = ref(false);
const userDialog = ref(false);
const deleteUserDialog = ref(false);
const submitted = ref(false);
const isNewUser = ref(false);
const toast = useToast();
const dt = ref();
const defaultUser = {
  id: 0,
  name: 'otto',
  email: 'o@o.com',
  username: 'ouser',
  password: '',
  password_hash: '',
  external_id: 0,
  orgunit_id: 0,
  groups_id: [],
  phone: null,
  comment: null,
  is_admin: false,
  is_locked: false,
  is_active: true,
};
const groupsList = ref([]);
const dataCurrentUser = ref(defaultUser);
const dataUsers = ref([defaultUser]);
const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
  username: { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const props = defineProps({
  display: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['user-invalid-session']);

const getNumUsers = computed(() => {
  const method = 'getNumUsers';
  log.t(`##-->${moduleName}::${method}()`);
  if ((props.display === undefined) || props.display === false) {
    return 0;
  }
  if (loadedData.value) {
    return dataUsers.value.length;
  } if ((isNullOrUndefined(dataUsers.value)) || dataUsers.value.length < 1 || !loadingData.value) {
    user.getList((retval, statusMessage) => {
      if (statusMessage === 'SUCCESS') {
        dataUsers.value = retval;
        log.l(`# IN loadData -> dataUsers.value.length : ${dataUsers.value.length}`);
        loadedData.value = true;
        loadingData.value = false;
        return dataUsers.value.length;
      }
      log.e(`# GOT ERROR calling user.getList : ${statusMessage}, \n error:`, retval);
      loadedData.value = false;
      loadingData.value = false;
      return 0;
    });
  }
  return 0;
});

const findIndexById = (id) => {
  const method = 'findIndexById';
  log.t(`##-->${moduleName}::${method}(${id})`);
  let index = -1;
  for (let i = 0; i < dataUsers.value.length; i += 1) {
    if (dataUsers.value[i].id === id) {
      index = i;
      break;
    }
  }
  return index;
};
const checkNetworkError = (err) => {
  if (!isNullOrUndefined(err.response)) {
    log.w('retval.response', err.response, err.response.status, (err.response.status === 401));
    if (err.response.status === 401) {
      emit('user-invalid-session', 'User session is invalid', err.response);
    }
  }
};

const openNew = () => {
  const method = 'openNew';
  log.t(`##-->${moduleName}::${method}`);
  dataCurrentUser.value = defaultUser;
  submitted.value = false;
  isNewUser.value = true;
  userDialog.value = true;
};
const hideDialog = () => {
  const method = 'hideDialog';
  log.t(`##-->${moduleName}::${method}`);
  userDialog.value = false;
  isNewUser.value = false;
  submitted.value = false;
};

const loadGroupsList = (fnOnsuccess = null) => {
  const method = 'loadGroupsList';
  log.t(`##-->${moduleName}::${method}()`);
  if ((isNullOrUndefined(groupsList.value)) || !loadingGroupsList.value) {
    loadingGroupsList.value = true;
    group.getList((retval, statusMessage) => {
      if (statusMessage === 'SUCCESS') {
        // const listOfGroups = [{ name: 'readers', value: '1' }]; // retval.map((e)=>{ e.})
        const listOfGroups = retval.map((e) => ({ id: e.id, name: e.name }));
        groupsList.value = listOfGroups;
        log.l(`# IN getListGroups -> dataGroups.value.length : ${groupsList.value.length}`);
        loadingGroupsList.value = false;
        fnOnsuccess();
        return;
      }
      log.e(`# GOT ERROR calling group.getList : ${statusMessage}, \n error:`, retval);
      loadingGroupsList.value = false;
    });
  }
};

const saveUser = () => {
  const method = 'saveUser';
  log.t(`##-->${moduleName}::${method} id:[${dataCurrentUser.value.id}`);
  submitted.value = true;

  if (dataCurrentUser.value.name.trim()) {
    if (dataCurrentUser.value.id) { // SAVE EDITED USER
      log.l(`##-->${moduleName}::${method} SAVING EDITED USER id:[${dataCurrentUser.value.id}]`);
      isNewUser.value = false;
      const tempUser = { ...dataCurrentUser.value };
      if (!isNullOrUndefined(tempUser.password) && tempUser.password.trim().length > 2) {
        tempUser.password_hash = `${getPasswordHash(tempUser.password)}`;
      } else {
        tempUser.password_hash = '';
      }
      user.modifyUser(tempUser, (retval, statusMessage) => {
        if (statusMessage === 'SUCCESS') {
          toast.add({
            severity: 'success', summary: 'Successful', detail: 'User Updated', life: timeToDisplaySucces,
          });
          log.l('# in save callback for user.modifyUser call val', retval);
          dataUsers.value[findIndexById(retval.id)] = retval;
          log.l('# findIndexById(retval.id) : ', retval.id, findIndexById(retval.id));
          log.l('# dataUsers.value : ', dataUsers.value);
          // this.initialize(); // on recupere la liste a jour
        } else {
          log.e(`##-->${moduleName}::${method} ERROR SAVING EDITED USER id:[${dataCurrentUser.value.id}] ERROR : ${statusMessage} \n error:`, retval);
          toast.add({
            severity: 'error', summary: 'Error', detail: `⚡⚡⚠ User was not saved in DB ! error: ${statusMessage}`, life: timeToDisplayError,
          });
          checkNetworkError(retval);
        }
      });
    } else { // NEW USER
      log.l(`##-->${moduleName}::${method} SAVING NEW USER`);
      isNewUser.value = true;
      const tempUser = { ...dataCurrentUser.value };
      tempUser.id = 0;
      tempUser.password_hash = `${getPasswordHash(dataCurrentUser.value.password)}`;
      user.newUser(tempUser, (retval, statusMessage) => {
        if (statusMessage === 'SUCCESS') {
          log.w('# in saveDialog callback for user.newUser call val', retval);
          toast.add({
            severity: 'success', summary: 'Successful', detail: `User created in DB id: ${retval.id}`, life: timeToDisplaySucces,
          });
          log.w(`# in saveDialog for new item id ${retval}`);
          tempUser.datecreated = new Date();
          tempUser.id = retval.id;
          dataUsers.value.push(tempUser);
          // this.initialize(); // on recupere la liste a jour
        } else {
          log.e(`# ERROR in saveDialog callback for objet.newObjet call ERROR : ${statusMessage} \n error:`, retval);
          toast.add({
            severity: 'error', summary: 'Error', detail: `⚡⚡⚠ User was not created in DB ! error: ${statusMessage}`, life: timeToDisplayError,
          });
          checkNetworkError(retval);
        }
      });
    }
    userDialog.value = false;
    dataCurrentUser.value = defaultUser;
  }
};
const editUser = (currentUser) => {
  const method = 'editUser';
  log.t(`##-->${moduleName}::${method}(id:${currentUser.id})`, currentUser);
  dataCurrentUser.value = { ...currentUser };
  user.getUser(currentUser.id, (retval, statusMessage) => {
    if (statusMessage === 'SUCCESS') {
      log.l('# in editUser callback for user.getUser call val', retval);
      dataCurrentUser.value = { ...retval };
      // let's retrieve all groups
      loadGroupsList(() => {
        userDialog.value = true;
      });
    } else {
      log.e(`# ERROR in editUser user.getUser callback: ${statusMessage} \n error:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to retrieve this User id: [${currentUser.id} from DB ! error: ${statusMessage}`, life: timeToDisplayError,
      });
      checkNetworkError(retval);
    }
  });
};
const confirmDeleteUser = (currentUser) => {
  const method = 'confirmDeleteUser';
  log.t(`##-->${moduleName}::${method}`);
  if (currentUser.id === getUserId()) {
    toast.add({
      severity: 'error', summary: 'Error', detail: `⚡⚡⚠ You cannot erase your own account ! User id: [${currentUser.id} `, life: timeToDisplayError,
    });
    dataCurrentUser.value = defaultUser;
    deleteUserDialog.value = false;
    return;
  }
  dataCurrentUser.value = currentUser;
  deleteUserDialog.value = true;
};
const deleteUser = () => {
  const method = 'deleteUser';
  const { id } = dataCurrentUser.value;
  log.t(`##-->${moduleName}::${method}(id:${id})`);
  if (id === getUserId()) {
    toast.add({
      severity: 'error', summary: 'Error', detail: `⚡⚡⚠ You cannot erase your own account ! User id: [${dataCurrentUser.value.id} `, life: timeToDisplayError,
    });
    return;
  }
  user.deleteUser(id, (retval, statusMessage) => {
    if (statusMessage === 'SUCCESS') {
      log.w('# in saveDialog callback for user.deleteUser call val', retval);
      dataUsers.value = dataUsers.value.filter((val) => val.id !== dataCurrentUser.value.id);
      deleteUserDialog.value = false;
      dataCurrentUser.value = defaultUser;
      toast.add({
        severity: 'success', summary: 'Successful', detail: 'User Deleted', life: timeToDisplaySucces,
      });
    } else {
      log.e(`# ERROR in deleteUser callback for user.deleteUser call ERROR : ${statusMessage} \n error:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to delete this User id: [${dataCurrentUser.value.id} from DB ! error: ${statusMessage}`, life: timeToDisplayError,
      });
      checkNetworkError(retval);
    }
  });
};

onMounted(() => {
  const method = 'onMounted';
  log.t(`##-->${moduleName}::${method}`);
});

</script>

<style>
.align-right {
  text-align: right !important;
}

.align-center {
  text-align: center;
}

.p-datatable th[class*="align-center"] .p-column-header-content {
  justify-content: center;
}

.p-datatable td[class*="align-right"] .p-column-header-content {
  justify-content: end;
  text-align: right;
}

.cg-centered-header {
  justify-content: center;
  align-content: center;
  background-color: blue;
}
.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.confirmation-content {
  display: flex;
  align-items: center;
  justify-content: center;
}

</style>
