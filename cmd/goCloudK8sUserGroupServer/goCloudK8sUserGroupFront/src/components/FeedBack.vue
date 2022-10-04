<template>
  <transition name="p-message" appear>
    <div v-show="feedbackVisible" :class="containerClass" role="alert">
      <div class="p-message-wrapper">
        <span :class="iconClass" />
        <div class="p-message-text">
          {{ feedbackText }}
        </div>
        <button v-if="closable" v-ripple class="p-message-close p-link" type="button" @click="feedbackVisible = false">
          <i class="p-message-close-icon pi pi-times" />
        </button>
      </div>
    </div>
  </transition>
</template>

<script setup>
import {
  computed, onMounted, ref, watch,
} from 'vue';
import vRipple from 'primevue/ripple';
import { getLog } from '../config';

const moduleName = 'FeedBack';
const validTypes = ['success', 'info', 'warn', 'error'];
const feedbackVisible = ref(false);
const feedbackText = ref('');
const feedbackType = ref('info');
const log = getLog(moduleName, 4, 2);
const props = defineProps({
  msg: {
    type: String,
    required: true,
    default: '',
  },
  msgType: {
    type: String,
    required: true,
    default: 'info',
  },
  visible: {
    type: Boolean,
    required: true,
    default: false,
  },
  closable: {
    type: Boolean,
    default: true,
  },
  icon: {
    type: String,
    default: null,
  },
});

watch(() => props.msg, () => {
  log.t(`watch props.msg : ${props.msg}`);
  feedbackText.value = props.msg;
});

watch(() => props.msgType, () => {
  log.t(`watch props.msgType : ${props.msgType}`);
  if (validTypes.includes(props.msgType)) {
    feedbackType.value = props.msgType;
  } else {
    feedbackType.value = 'info';
  }
});

watch(() => props.visible, () => {
  log.t(`watch props.displayIt : ${props.visible}`);
  feedbackVisible.value = props.visible;
});

const containerClass = computed(() => `p-message p-component p-message-${feedbackType.value}`);
const iconClass = computed(() => [
  'p-message-icon pi',
  props.icon
    ? props.icon
    : {
      'pi-info-circle': feedbackType.value === 'info',
      'pi-check': feedbackType.value === 'success',
      'pi-exclamation-triangle': feedbackType.value === 'warn',
      'pi-times-circle': feedbackType.value === 'error',
    },
]);

const displayFeedBack = (text, type) => {
  log.t(`displayFeedBack() text:'${text}' type:'${type}'`);
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

defineExpose({ displayFeedBack, resetFeedBack });
onMounted(() => {
  log.t(`mounted() msg:'${props.msg}' type:'${props.msgType}'`);
  feedbackText.value = props.msg;
  feedbackType.value = props.msgType;
  feedbackVisible.value = props.visible;
});
</script>

<style>
.p-message-wrapper {
  display: flex;
  align-items: center;
}

.p-message-close {
  display: flex;
  align-items: center;
  justify-content: center;
}

.p-message-close.p-link {
  margin-left: auto;
  overflow: hidden;
  position: relative;
}

.p-message-enter-from {
  opacity: 0;
}

.p-message-enter-active {
  transition: opacity 0.3s;
}

.p-message.p-message-leave-from {
  max-height: 1000px;
}

.p-message.p-message-leave-to {
  max-height: 0;
  opacity: 0;
  margin: 0 !important;
}

.p-message-leave-active {
  overflow: hidden;
  transition: max-height 0.3s cubic-bezier(0, 1, 0, 1), opacity 0.3s, margin 0.15s;
}

.p-message-leave-active .p-message-close {
  display: none;
}
</style>
