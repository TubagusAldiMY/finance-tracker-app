<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  label: String,
  type: {
    type: String,
    default: 'text'
  },
  placeholder: String,

  error: {
    type: String,
    default: ''
  },

  inputClass: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const baseClass =
  'w-full px-4 py-2 border rounded-lg focus:ring-2 transition'

const errorClass =
  'border-red-500 focus:ring-red-400'

const normalClass =
  'border-slate-300 focus:ring-blue-500'

const computedClass = computed(() => {
  return [
    baseClass,
    props.error ? errorClass : normalClass,
    props.inputClass
  ].join(' ')
})
</script>

<template>
  <div class="space-y-1">
    <label v-if="label" class="block text-sm font-medium text-slate-700">
      {{ label }}
    </label>

    <input
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :class="computedClass"
      @input="emit('update:modelValue', $event.target.value)"
    />

    <p
      v-if="error"
      class="text-sm text-red-500"
    >
      {{ error }}
    </p>
  </div>
</template>
