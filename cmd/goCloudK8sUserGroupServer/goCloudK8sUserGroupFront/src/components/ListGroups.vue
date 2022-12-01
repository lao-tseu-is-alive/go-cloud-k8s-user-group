<template>
  <div class="w-full card">
    <DataTable
      ref="dt"
      v-model:filters="filters"
      :value="dataGroups"
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
              Liste des {{ getNumGroups }} groupes
            </h4>
          </template>
          <template #end>
            <Button label="Nouvel Groupe" icon="pi pi-plus" class="p-button-success mr-0" @click="openNew" />
          </template>
        </Toolbar>
      </template>
      <ColumnGroup type="header">
        <Row>
          <Column field="id" header="id" :sortable="true" />
          <Column field="name" header="Nom" :sortable="true" />
          <Column field="is_active" header="Active?" :sortable="true" />
          <Column header="Actions" />
        </Row>
      </ColumnGroup>
      <Column field="id" />
      <Column field="name" class="align-right" />
      <Column field="is_active" class="align-center" />
      <Column :exportable="false" style="min-width:8rem">
        <template #body="slotProps">
          <Button icon="pi pi-pencil" class="p-button-rounded  mr-2" @click="editGroup(slotProps.data)" />
          <Button icon="pi pi-trash" class="p-button-rounded " @click="confirmDeleteGroup(slotProps.data)" />
        </template>
      </Column>
    </DataTable>
  </div>
  <!-- BEGIN DIALOG EDIT USER -->
  <Dialog v-model:visible="groupDialog" :style="{ width: '450px' }" :header="`Group id [${dataCurrentGroup.id}] details`" :modal="true" class="p-fluid">
    <div class="field">
      <label for="name">Name</label>
      <InputText id="name" v-model.trim="dataCurrentGroup.name" required="true" autofocus :class="{ 'p-invalid': submitted && !dataCurrentGroup.name }" />
      <small v-if="submitted && !dataCurrentGroup.name" class="p-error">Name is required.</small>
    </div>
    <div class="field">
      <label for="description">Comment</label>
      <Textarea id="description" v-model="dataCurrentGroup.comment" required="true" rows="3" cols="20" />
    </div>

    <div class="formgrid grid">
      <div class="field col">
        <label for="isactive">IS active ?</label>
        <InputSwitch id="isactive" v-model="dataCurrentGroup.is_active" />
      </div>
    </div>
    <template #footer>
      <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="hideDialog" />
      <Button label="Save" icon="pi pi-check" class="p-button-text" @click="saveGroup" />
    </template>
  </Dialog>
  <!-- END DIALOG EDIT USER -->
  <Dialog v-model:visible="deleteGroupDialog" :style="{ width: '450px' }" header="Confirm" :modal="true">
    <div class="confirmation-content">
      <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
      <span v-if="dataCurrentGroup">
        Are you sure you want to delete <b>{{ dataCurrentGroup.name }}</b>?
      </span>
    </div>
    <template #footer>
      <Button label="No" icon="pi pi-times" class="p-button-text" @click="deleteGroupDialog = false" />
      <Button label="Yes" icon="pi pi-check" class="p-button-text" @click="deleteGroup" />
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
import group from './Group';
import { getLog } from '../config';
import { isNullOrUndefined } from '../tools/utils';

const moduleName = 'ListGroups';
const timeToDisplayError = 7000;
const timeToDisplaySucces = 4000;

const log = getLog(moduleName, 4, 2);
const loadedData = ref(false);
const loadingData = ref(false);
const groupDialog = ref(false);
const deleteGroupDialog = ref(false);
const submitted = ref(false);
const isNewGroup = ref(false);
const toast = useToast();
const dt = ref();
const defaultGroup = {
  id: 0,
  name: 'nouveau groupe',
  comment: null,
  is_active: true,
};
const dataCurrentGroup = ref(defaultGroup);
const dataGroups = ref([defaultGroup]);
const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
  name: { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const props = defineProps({
  display: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['user-invalid-session']);

const getNumGroups = computed(() => {
  const method = 'getNumGroups';
  log.t(`##-->${moduleName}::${method}()`);
  if ((props.display === undefined) || props.display === false) {
    return 0;
  }
  if (loadedData.value) {
    return dataGroups.value.length;
  }
  if ((isNullOrUndefined(dataGroups.value)) || dataGroups.value.length < 1 || !loadingData.value) {
    group.getList((retval, statusMessage) => {
      if (statusMessage === 'SUCCESS') {
        dataGroups.value = retval;
        log.l(`# IN loadData -> dataGroups.value.length : ${dataGroups.value.length}`);
        loadedData.value = true;
        loadingData.value = false;
        return dataGroups.value.length;
      }
      log.e(`# GOT ERROR calling group.getList : ${statusMessage}, \n error:`, retval);
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
  for (let i = 0; i < dataGroups.value.length; i += 1) {
    if (dataGroups.value[i].id === id) {
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
  dataCurrentGroup.value = defaultGroup;
  submitted.value = false;
  isNewGroup.value = true;
  groupDialog.value = true;
};
const hideDialog = () => {
  const method = 'hideDialog';
  log.t(`##-->${moduleName}::${method}`);
  groupDialog.value = false;
  isNewGroup.value = false;
  submitted.value = false;
};
const saveGroup = () => {
  const method = 'saveGroup';
  log.t(`##-->${moduleName}::${method} id:[${dataCurrentGroup.value.id}`);
  submitted.value = true;

  if (dataCurrentGroup.value.name.trim()) {
    if (dataCurrentGroup.value.id) { // SAVE EDITED USER
      log.l(`##-->${moduleName}::${method} SAVING EDITED USER id:[${dataCurrentGroup.value.id}]`);
      isNewGroup.value = false;
      const tempGroup = { ...dataCurrentGroup.value };
      group.modifyGroup(tempGroup, (retval, statusMessage) => {
        if (statusMessage === 'SUCCESS') {
          toast.add({
            severity: 'success', summary: 'Successful', detail: 'Group Updated', life: timeToDisplaySucces,
          });
          log.l('# in save callback for group.modifyGroup call val', retval);
          dataGroups.value[findIndexById(retval.id)] = retval;
          log.l('# findIndexById(retval.id) : ', retval.id, findIndexById(retval.id));
          log.l('# dataGroups.value : ', dataGroups.value);
          // this.initialize(); // on recupere la liste a jour
        } else {
          log.e(`##-->${moduleName}::${method} ERROR SAVING EDITED USER id:[${dataCurrentGroup.value.id}] ERROR : ${statusMessage} \n error:`, retval);
          toast.add({
            severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Group was not saved in DB ! error: ${statusMessage}`, life: timeToDisplayError,
          });
          checkNetworkError(retval);
        }
      });
    } else { // NEW USER
      log.l(`##-->${moduleName}::${method} SAVING NEW USER`);
      isNewGroup.value = true;
      const tempGroup = { ...dataCurrentGroup.value };
      tempGroup.id = 0;
      group.newGroup(tempGroup, (retval, statusMessage) => {
        if (statusMessage === 'SUCCESS') {
          log.w('# in saveDialog callback for group.newGroup call val', retval);
          toast.add({
            severity: 'success', summary: 'Successful', detail: `Group created in DB id: ${retval.id}`, life: timeToDisplaySucces,
          });
          log.w(`# in saveDialog for new item id ${retval}`);
          tempGroup.datecreated = new Date();
          tempGroup.id = retval.id;
          dataGroups.value.push(tempGroup);
          // this.initialize(); // on recupere la liste a jour
        } else {
          log.e(`# ERROR in saveDialog callback for objet.newObjet call ERROR : ${statusMessage} \n error:`, retval);
          toast.add({
            severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Group was not created in DB ! error: ${statusMessage}`, life: timeToDisplayError,
          });
          checkNetworkError(retval);
        }
      });
    }
    groupDialog.value = false;
    dataCurrentGroup.value = defaultGroup;
  }
};
const editGroup = (currentGroup) => {
  const method = 'editGroup';
  log.t(`##-->${moduleName}::${method}`, currentGroup);
  dataCurrentGroup.value = { ...currentGroup };
  group.getGroup(currentGroup.id, (retval, statusMessage) => {
    if (statusMessage === 'SUCCESS') {
      log.l('# in editGroup callback for group.getGroup call val', retval);
      dataCurrentGroup.value = { ...retval };
      groupDialog.value = true;
    } else {
      log.e(`# ERROR in editGroup group.getGroup callback: ${statusMessage} \n error:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to retrieve this Group id: [${currentGroup.id} from DB ! error: ${statusMessage}`, life: timeToDisplayError,
      });
      checkNetworkError(retval);
    }
  });
};
const confirmDeleteGroup = (currentGroup) => {
  const method = 'confirmDeleteGroup';
  log.t(`##-->${moduleName}::${method}`);
  dataCurrentGroup.value = currentGroup;
  deleteGroupDialog.value = true;
};
const deleteGroup = () => {
  const method = 'deleteGroup';
  const { id } = dataCurrentGroup.value;
  log.t(`##-->${moduleName}::${method}(id:${id})`);
  group.deleteGroup(id, (retval, statusMessage) => {
    if (statusMessage === 'SUCCESS') {
      log.w('# in saveDialog callback for group.deleteGroup call val', retval);
      dataGroups.value = dataGroups.value.filter((val) => val.id !== dataCurrentGroup.value.id);
      deleteGroupDialog.value = false;
      dataCurrentGroup.value = defaultGroup;
      toast.add({
        severity: 'success', summary: 'Successful', detail: 'Group Deleted', life: timeToDisplaySucces,
      });
    } else {
      log.e(`# ERROR in deleteGroup callback for group.deleteGroup call ERROR : ${statusMessage} \n error:`, retval);
      toast.add({
        severity: 'error', summary: 'Error', detail: `⚡⚡⚠ Unable to delete this Group id: [${dataCurrentGroup.value.id} from DB ! error: ${statusMessage}`, life: timeToDisplayError,
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
