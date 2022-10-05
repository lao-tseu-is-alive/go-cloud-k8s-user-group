import { createApp } from 'vue';
import PrimeVue from 'primevue/config';
import ToastService from 'primevue/toastservice';
import App from './App.vue';
import 'primeflex/primeflex.css';
import 'primevue/resources/themes/lara-light-blue/theme.css';
import 'primevue/resources/primevue.min.css'; // core css
import 'primeicons/primeicons.css';

// https://www.primefaces.org/primevue/setup
// free themes to choose from :
/*
 primevue/resources/themes/bootstrap4-light-blue/theme.css
 primevue/resources/themes/bootstrap4-light-purple/theme.css
 primevue/resources/themes/bootstrap4-dark-blue/theme.css
 primevue/resources/themes/bootstrap4-dark-purple/theme.css
 primevue/resources/themes/md-light-indigo/theme.css
 primevue/resources/themes/md-light-deeppurple/theme.css
 primevue/resources/themes/md-dark-indigo/theme.css
 primevue/resources/themes/md-dark-deeppurple/theme.css
 primevue/resources/themes/mdc-light-indigo/theme.css
 primevue/resources/themes/mdc-light-deeppurple/theme.css
 primevue/resources/themes/mdc-dark-indigo/theme.css
 primevue/resources/themes/mdc-dark-deeppurple/theme.css
 primevue/resources/themes/tailwind-light/theme.css
 primevue/resources/themes/fluent-light/theme.css
 primevue/resources/themes/lara-light-indigo/theme.css
 primevue/resources/themes/lara-dark-indigo/theme.css
 primevue/resources/themes/lara-light-purple/theme.css
 primevue/resources/themes/lara-dark-purple/theme.css
 primevue/resources/themes/lara-light-blue/theme.css
 primevue/resources/themes/lara-dark-blue/theme.css
 primevue/resources/themes/lara-light-teal/theme.css
 primevue/resources/themes/lara-dark-teal/theme.css
 primevue/resources/themes/saga-blue/theme.css
 primevue/resources/themes/saga-green/theme.css
 primevue/resources/themes/saga-orange/theme.css
 primevue/resources/themes/saga-purple/theme.css
 primevue/resources/themes/vela-blue/theme.css
 primevue/resources/themes/vela-green/theme.css
 primevue/resources/themes/vela-orange/theme.css
 primevue/resources/themes/vela-purple/theme.css
 primevue/resources/themes/arya-blue/theme.css
 primevue/resources/themes/arya-green/theme.css
 primevue/resources/themes/arya-orange/theme.css
 primevue/resources/themes/arya-purple/theme.css
 primevue/resources/themes/nova/theme.css
 primevue/resources/themes/nova-alt/theme.css
 primevue/resources/themes/nova-accent/theme.css
 primevue/resources/themes/nova-vue/theme.css
 primevue/resources/themes/luna-amber/theme.css
 primevue/resources/themes/luna-blue/theme.css
 primevue/resources/themes/luna-green/theme.css
 primevue/resources/themes/luna-pink/theme.css
 primevue/resources/themes/rhea/theme.css
 */

const app = createApp(App);
app.use(PrimeVue);
app.use(ToastService);
app.mount('#app');
