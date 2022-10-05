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
  <Dialog v-model:visible="userDialog" :style="{ width: '450px' }" header="User Details" :modal="true" class="p-fluid">
    <div class="field">
      <label for="name">Name</label>
      <InputText id="name" v-model.trim="dataNewUser.name" required="true" autofocus :class="{ 'p-invalid': submitted && !dataNewUser.name }" />
      <small v-if="submitted && !dataNewUser.name" class="p-error">Name is required.</small>
    </div>
    <div class="field">
      <label for="name">Username</label>
      <InputText id="username" v-model.trim="dataNewUser.username" required="true" :class="{ 'p-invalid': submitted && !dataNewUser.username }" />
      <small v-if="submitted && !dataNewUser.username" class="p-error">Username is required.</small>
    </div>
    <div class="field">
      <label for="name">Password</label>
      <InputText id="email" v-model.trim="dataNewUser.password_hash" required="true" :class="{ 'p-invalid': submitted && !dataNewUser.password_hash }" />
      <small v-if="submitted && !dataNewUser.password_hash" class="p-error">Password is required.</small>
    </div>
    <div class="field">
      <label for="name">Email</label>
      <InputText id="email" v-model.trim="dataNewUser.email" required="true" :class="{ 'p-invalid': submitted && !dataNewUser.email }" />
      <small v-if="submitted && !dataNewUser.email" class="p-error">Email is required.</small>
    </div>
    <div class="field">
      <label for="name">OrgUnit</label>
      <InputText id="email" v-model.trim="dataNewUser.enterprise" />
    </div>
    <div class="field">
      <label for="name">Phone</label>
      <InputText id="email" v-model.trim="dataNewUser.phone" />
    </div>
    <div class="field">
      <label for="description">Comment</label>
      <Textarea id="description" v-model="dataNewUser.comment" required="true" rows="3" cols="20" />
    </div>

    <div class="formgrid grid">
      <div class="field col">
        <label for="isadmin">Is Administrator ?</label>
        <InputSwitch id="isadmin" v-model="dataNewUser.is_admin" />
      </div>
      <div class="field col">
        <label for="islocked">IS locked ?</label>
        <InputSwitch id="islocked" v-model="dataNewUser.is_locked" />
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
      <span v-if="dataNewUser">Are you sure you want to delete <b>{{ dataNewUser.name }}</b>?</span>
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
import { getLog } from '../config';
import { isNullOrUndefined } from '../tools/utils';

const moduleName = 'ListUsers';

const log = getLog(moduleName, 4, 2);
const loadedData = ref(false);
const userDialog = ref(false);
const deleteUserDialog = ref(false);
const submitted = ref(false);
const toast = useToast();
const dt = ref();
// u.Name, u.Email, u.Username, u.PasswordHash, &u.ExternalId, &u.Enterprise, &u.Phone, //$1-$7
// u.IsAdmin, u.Creator, &u.Comment
const defaultUser = {
  id: 0,
  name: '',
  email: '',
  username: '',
  password_hash: '',
  enterprise: null,
  phone: null,
  comment: null,
  is_admin: false,
  is_locked: false,
  is_active: true,
};
const dataNewUser = ref(defaultUser);
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
  if ((props.display === undefined) || props.display === false) {
    return 0;
  }
  if ((isNullOrUndefined(dataUsers.value)) || dataUsers.value.length < 2) {
    user.getList((data, errMessage) => {
      if (!isNullOrUndefined(data)) {
        dataUsers.value = data;
        log.l('# IN loadData -> objet data :', dataUsers.value);
        loadedData.value = true;
        return dataUsers.value.length;
      }
      log.e(`# GOT ERROR IN loadData() in callback for user.getList data: ${errMessage}`);
      loadedData.value = false;
      return 0;
    });
  }
  return dataUsers.value.length;
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
  dataNewUser.value = {
    id: 0,
    username: '',
    name: '',
    is_admin: false,
    is_locked: false,
  };
  submitted.value = false;
  userDialog.value = true;
};
const hideDialog = () => {
  const method = 'hideDialog';
  log.t(`##-->${moduleName}::${method}`);
  userDialog.value = false;
  submitted.value = false;
};
const saveUser = () => {
  const method = 'saveUser';
  log.t(`##-->${moduleName}::${method}`);
  submitted.value = true;

  if (dataNewUser.value.name.trim()) {
    if (dataNewUser.value.id) {
      dataUsers.value[findIndexById(dataNewUser.value.id)] = dataNewUser.value;
      // ('⚡⚡⚠ PAS DE RESEAU ! ☹ vous êtes "OFFLINE" ', 'error');
      toast.add({
        severity: 'success', summary: 'Successful', detail: 'User Updated', life: 3000,
      });
    } else {
      dataNewUser.value.id = 9999;
      dataUsers.value.push(dataNewUser.value);
      toast.add({
        severity: 'success', summary: 'Successful', detail: 'User Created', life: 3000,
      });
    }

    userDialog.value = false;
    dataNewUser.value = {};
  }
};
const editUser = (currentUser) => {
  const method = 'editUser';
  log.t(`##-->${moduleName}::${method}`, currentUser);
  dataNewUser.value = { ...currentUser };
  userDialog.value = true;
};
const confirmDeleteUser = (currentUser) => {
  const method = 'confirmDeleteUser';
  log.t(`##-->${moduleName}::${method}`);
  dataNewUser.value = currentUser;
  deleteUserDialog.value = true;
};
const deleteUser = () => {
  const method = 'deleteUser';
  log.t(`##-->${moduleName}::${method}`);
  dataUsers.value = dataUsers.value.filter((val) => val.id !== dataNewUser.value.id);
  deleteUserDialog.value = false;
  dataNewUser.value = {};
  toast.add({
    severity: 'success', summary: 'Successful', detail: 'User Deleted', life: 3000,
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
