<template>
  <div class="w-full">
    <Card>
      <template #title>
        <h4 class="m-2">
          Liste des {{ getNumUsers }} utilisateurs
        </h4>
      </template>
      <template #content>
        <DataTable
          v-model:filters="filters"
          :value="dataUsers"
          removable-sort
          responsive-layout="scroll"
        >
          <template #header>
            <div class="flex align-items-center justify-content-between">
              <span class="p-input-icon-left right-0">
                <i class="pi pi-search" />
                <InputText
                  v-model="filters['username'].value"
                  placeholder="filtre de couche"
                />
              </span>
            </div>
          </template>
          <ColumnGroup type="header">
            <Row>
              <Column field="id" header="id" :sortable="true" />
              <Column field="username" header="Username" :sortable="true" />
              <Column field="name" header="Nom" :sortable="true" />
              <Column field="is_admin" header="Admin?" :sortable="true" />
              <Column field="is_locked" header="Locked?" :sortable="true" />
            </Row>
          </ColumnGroup>
          <Column field="id" />
          <Column field="username" class="align-right" />
          <Column field="name" class="align-right" />
          <Column field="is_admin" class="align-right" />
          <Column field="is_locked" class="align-right" />
        </DataTable>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue';
import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import ColumnGroup from 'primevue/columngroup';
import Row from 'primevue/row';
import InputText from 'primevue/inputtext';
import { FilterMatchMode } from 'primevue/api';
import user from './User';
import { getLog } from '../config';
import { isNullOrUndefined } from '../tools/utils';

const moduleName = 'ListUsers';

const log = getLog(moduleName, 4, 2);

const loadedData = ref(false);
const dataUsers = ref([
  {
    id: 0,
    username: 'aucune données n\'est disponible en ce moment...',
    name: '',
    is_admin: false,
    is_locked: false,
  },
]);
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

onMounted(() => {
  const method = 'onMounted';
  log.l(`##-->${moduleName}::${method}`);
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
</style>
