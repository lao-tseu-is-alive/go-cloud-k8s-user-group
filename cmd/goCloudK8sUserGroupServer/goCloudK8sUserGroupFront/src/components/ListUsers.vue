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
            <Button label="New User" icon="pi pi-plus" class="p-button-success mr-0" @click="openNew" />
          </template>
        </Toolbar>
      </template>
      <ColumnGroup type="header">
        <Row>
          <Column field="id" header="id" :sortable="true" />
          <Column field="username" header="Username" :sortable="true" />
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
  <Dialog v-model:visible="userDialog" :style="{ width: '450px' }" :header="`User id [${dataCurrentUser.id}] details`" :modal="true" class="p-fluid">
    <div class="field">
      <label for="name">Name</label>
      <InputText id="name" v-model.trim="dataCurrentUser.name" required="true" autofocus :class="{ 'p-invalid': submitted && !dataCurrentUser.name }" />
      <small v-if="submitted && !dataCurrentUser.name" class="p-error">Name is required.</small>
    </div>
    <div class="field">
      <label for="name">Username</label>
      <InputText id="username" v-model.trim="dataCurrentUser.username" required="true" :class="{ 'p-invalid': submitted && !dataCurrentUser.username }" />
      <small v-if="submitted && !dataCurrentUser.username" class="p-error">Username is required.</small>
    </div>
    <div class="field">
      <label for="name">Password</label>
      <InputText id="password" v-model.trim="dataCurrentUser.password" type="password" :class="{ 'p-invalid': submitted && isNewUser && !dataCurrentUser.password }" />
      <small v-if="submitted && isNewUser && !dataCurrentUser.password" class="p-error">Password is required.</small>
    </div>
    <div class="field">
      <label for="name">Email</label>
      <InputText id="email" v-model.trim="dataCurrentUser.email" required="true" :class="{ 'p-invalid': submitted && !dataCurrentUser.email }" />
      <small v-if="submitted && !dataCurrentUser.email" class="p-error">Email is required.</small>
    </div>
    <div class="field">
      <label for="name">OrgUnit</label>
      <InputText id="email" v-model.trim="dataCurrentUser.enterprise" />
    </div>
    <div class="field">
      <label for="name">Phone</label>
      <InputText id="email" v-model.trim="dataCurrentUser.phone" />
    </div>
    <div class="field">
      <label for="description">Comment</label>
      <Textarea id="description" v-model="dataCurrentUser.comment" required="true" rows="3" cols="20" />
    </div>

    <div class="formgrid grid">
      <div class="field col">
        <label for="isadmin">Is Administrator ?</label>
        <InputSwitch id="isadmin" v-model="dataCurrentUser.is_admin" />
      </div>
      <div class="field col">
        <label for="islocked">IS locked ?</label>
        <InputSwitch id="islocked" v-model="dataCurrentUser.is_locked" />
      </div>
      <div class="field col">
        <label for="isactive">IS active ?</label>
        <InputSwitch id="isactive" v-model="dataCurrentUser.is_active" />
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
import Toolbar from 'primevue/toolbar';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import ColumnGroup from 'primevue/columngroup';
import Row from 'primevue/row';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import InputSwitch from 'primevue/inputswitch';
import Textarea from 'primevue/textarea';
import { FilterMatchMode } from 'primevue/api';
import user from './User';
import { getPasswordHash, getUserId } from './Login';
import { getLog } from '../config';
import { isNullOrUndefined } from '../tools/utils';

const moduleName = 'ListUsers';
const timeToDisplayError = 7000;
const timeToDisplaySucces = 4000;

const log = getLog(moduleName, 4, 2);
const loadedData = ref(false);
const loadingData = ref(false);
const userDialog = ref(false);
const deleteUserDialog = ref(false);
const submitted = ref(false);
const isNewUser = ref(false);
const toast = useToast();
const dt = ref();
// u.Name, u.Email, u.Username, u.PasswordHash, &u.ExternalId, &u.Enterprise, &u.Phone, //$1-$7
// u.IsAdmin, u.Creator, &u.Comment
const defaultUser = {
  id: 0,
  name: 'otto',
  email: 'o@o.com',
  username: 'ouser',
  password: '',
  password_hash: '',
  external_id: 0,
  enterprise: null,
  phone: null,
  comment: null,
  is_admin: false,
  is_locked: false,
  is_active: true,
};
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

const getNumUsers = computed(() => {
  const method = 'getNumUsers';
  log.t(`##-->${moduleName}::${method}()`);
  if ((props.display === undefined) || props.display === false) {
    return 0;
  }
  if (loadedData.value) {
    return dataUsers.value.length;
  } if ((isNullOrUndefined(dataUsers.value)) || dataUsers.value.length < 1 || !loadingData.value) {
    user.getList((data, errMessage) => {
      if (!isNullOrUndefined(data)) {
        dataUsers.value = data;
        log.l(`# IN loadData -> dataUsers.value.length : ${dataUsers.value.length}`);
        loadedData.value = true;
        loadingData.value = false;
        return dataUsers.value.length;
      }
      log.e(`# GOT ERROR IN loadData() in callback for user.getList data: ${errMessage}`);
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
      user.modifyUser(tempUser, (retval, errorMsg) => {
        if (errorMsg === 'SUCCESS') {
          toast.add({
            severity: 'success', summary: 'Successful', detail: 'User Updated', life: timeToDisplaySucces,
          });
          log.l('# in save callback for user.modifyUser call val', retval);
          dataUsers.value[findIndexById(retval.id)] = retval;
          log.l('# findIndexById(retval.id) : ', retval.id, findIndexById(retval.id));
          log.l('# dataUsers.value : ', dataUsers.value);
          // this.initialize(); // on recupere la liste a jour
        } else {
          log.e(`##-->${moduleName}::${method} ERROR SAVING EDITED USER id:[${dataCurrentUser.value.id}] ERROR : ${errorMsg}`, retval);
          toast.add({
            severity: 'error', summary: 'Error', detail: `⚡⚡⚠ User was not saved in DB ! error: ${errorMsg}`, life: timeToDisplayError,
          });
        }
      });
    } else { // NEW USER
      log.l(`##-->${moduleName}::${method} SAVING NEW USER`);
      isNewUser.value = true;
      const tempUser = { ...dataCurrentUser.value };
      tempUser.id = 0;
      tempUser.password_hash = `${getPasswordHash(dataCurrentUser.value.password)}`;
      user.newUser(tempUser, (retval, errorMsg) => {
        if (errorMsg === 'SUCCESS') {
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
          log.e(`# ERROR in saveDialog callback for objet.newObjet call ERROR : ${errorMsg} val`, retval);
          toast.add({
            severity: 'error', summary: 'Error', detail: `⚡⚡⚠ User was not created in DB ! error: ${errorMsg}`, life: timeToDisplayError,
          });
        }
      });
    }
    userDialog.value = false;
    dataCurrentUser.value = defaultUser;
  }
};
const editUser = (currentUser) => {
  const method = 'editUser';
  log.t(`##-->${moduleName}::${method}`, currentUser);
  dataCurrentUser.value = { ...currentUser };
  user.getUser(currentUser.id, (retval, errorMsg) => {
    if (errorMsg === 'SUCCESS') {
      log.l('# in editUser callback for user.getUser call val', retval);
      dataCurrentUser.value = { ...retval };
      userDialog.value = true;
    } else {
      log.e(`# ERROR in editItem callback for user.getUser call ERROR : ${errorMsg} val:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to retrieve this User id: [${currentUser.value.id} from DB ! error: ${errorMsg}`, life: timeToDisplayError,
      });
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
  user.deleteUser(id, (retval, errorMsg) => {
    if (errorMsg === 'SUCCESS') {
      log.w('# in saveDialog callback for user.deleteUser call val', retval);
      dataUsers.value = dataUsers.value.filter((val) => val.id !== dataCurrentUser.value.id);
      deleteUserDialog.value = false;
      dataCurrentUser.value = defaultUser;
      toast.add({
        severity: 'success', summary: 'Successful', detail: 'User Deleted', life: timeToDisplaySucces,
      });
    } else {
      log.e(`# ERROR in deleteUser callback for user.deleteUser call ERROR : ${errorMsg} val:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to delete this User id: [${dataCurrentUser.value.id} from DB ! error: ${errorMsg}`, life: timeToDisplayError,
      });
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
