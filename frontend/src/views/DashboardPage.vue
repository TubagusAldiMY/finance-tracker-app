<script setup>
import { ref, computed } from 'vue'

const currentDate = ref(new Date())

const endOfMonth = computed(() => {
  const d = currentDate.value
  return new Date(d.getFullYear(), d.getMonth() + 1, 0).getDate()
})
const startDayOfMonth = computed(() => {
  const d = currentDate.value
  const firstDay = new Date(d.getFullYear(), d.getMonth(), 1).getDay()
  return (firstDay + 6) % 7
})

const daysInMonth = computed(() => [
  ...Array(startDayOfMonth.value).fill(null),
  ...Array.from({ length: endOfMonth.value }, (_, i) => i + 1),
])

const weekDays = ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min']

function prevMonth() {
  const d = currentDate.value
  currentDate.value = new Date(d.getFullYear(), d.getMonth() - 1, 1)
}

function nextMonth() {
  const d = currentDate.value
  currentDate.value = new Date(d.getFullYear(), d.getMonth() + 1, 1)
}
</script>

<template>
  <div class="max-w-md mx-auto p-6 bg-white rounded-xl shadow-md">
    <h1 class="text-xl font-semibold text-gray-800 mb-4 text-center">
      ini adalah dashboard
    </h1>

    <div class="flex items-center justify-between mb-4">
      <button @click="prevMonth" class="px-3 py-1 rounded-md bg-gray-200 hover:bg-gray-300 transition">
        ←
      </button>

      <p class="font-medium text-gray-700">
        {{ currentDate.toLocaleString('id-ID', { month: 'long', year: 'numeric' }) }}
      </p>

      <button @click="nextMonth" class="px-3 py-1 rounded-md bg-gray-200 hover:bg-gray-300 transition">
        →
      </button>
    </div>
    <div class="grid grid-cols-7 gap-2 text-center text-sm font-medium text-gray-500 mb-2">
      <div v-for="day in weekDays" :key="day">
        {{ day }}
      </div>
    </div>
    <div class="grid grid-cols-7 gap-2 text-center">
      <div v-for="day in daysInMonth" :key="day" class="py-2 rounded-md border border-gray-200
               hover:bg-blue-100 cursor-pointer
               text-gray-700">
        {{ day }}
      </div>
    </div>
  </div>
</template>
