<template>
  <div class="w-full">
    <Splitter>
      <SplitterPanel>
        <Card>
          <template #title>
            <h4 class="m-2">
              {{ msg }}
            </h4>
          </template>
          <template #content>
            <DataTable
              v-model:filters="filters"
              :value="dataLayers"
              removable-sort
              responsive-layout="scroll"
            >
              <template #header>
                <div class="flex align-items-center justify-content-between">
                  <span class="p-input-icon-left right-0">
                    <i class="pi pi-search" />
                    <InputText
                      v-model="filters['layer'].value"
                      placeholder="filtre de couche"
                    />
                  </span>
                </div>
              </template>
              <ColumnGroup type="header">
                <Row>
                  <Column field="layer" header="Couche" :colspan="1" :rowspan="2" :sortable="true" />
                  <Column :colspan="2" class="align-center" header="GC interne" />
                  <Column header="GC public" :colspan="2" class="align-center" />
                </Row>
                <Row>
                  <Column field="requests_gc_internal" header="hits" :sortable="true" />
                  <Column field="visits_gc_internal" header="visitors" :sortable="true" />
                  <Column field="requests_gc_public" header="hits" :sortable="true" />
                  <Column field="visits_gc_public" header="visitors" :sortable="true" />
                </Row>
              </ColumnGroup>
              <Column field="layer" />
              <Column field="requests_gc_internal" class="align-right" />
              <Column field="visits_gc_internal" class="align-right" />
              <Column field="requests_gc_public" class="align-right" />
              <Column field="visits_gc_public" class="align-right" />
            </DataTable>
          </template>
        </Card>
      </SplitterPanel>
      <SplitterPanel>
        <Card>
          <template #content>
            <ChartLayers :data-chart="dataLayers" :filter-label="filters['layer'].value" />
          </template>
        </Card>
      </SplitterPanel>
    </Splitter>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import Card from 'primevue/card';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import ColumnGroup from 'primevue/columngroup';
import Row from 'primevue/row';
import InputText from 'primevue/inputtext';
import { FilterMatchMode } from 'primevue/api';
import { fetchData } from './fetchData';
import { BACKEND_URL } from '../config';

const moduleName = 'ListUsers';
const urlUsers = `${BACKEND_URL}/data/layers_stats.json`;
const loadedData = ref(false);
const errorDataFetch = ref(null);
const dataUsers = ref([
  {
    username: "aucune données n'est disponible en ce moment...",
  },
]);
const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
  layer: { value: null, matchMode: FilterMatchMode.CONTAINS },
});

onMounted(() => {
  const method = 'onMounted';
  console.log(`##-->${moduleName}::${method}`);
  fetchData(urlLayerStats, dataLayers, loadedData, errorDataFetch)
    .then(() => console.log(`##-->${moduleName}::${method} fetchData OK loadedData:${loadedData.value}`));
});

defineProps({
  msg: {
    type: String,
    required: true,
  },
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
