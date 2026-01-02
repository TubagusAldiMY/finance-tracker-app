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

const weeksInMonth = computed(() => {
  const weeks = []
  const days = daysInMonth.value

  for (let i = 0; i < days.length; i += 7) {
    weeks.push({
      index: i / 7,
      days: days.slice(i, i + 7),
    })
  }

  return weeks
})
const dailyBudgets = computed(() => Array.from({ length: endOfMonth.value }, () => 1000)) // Placeholder budget per day

const dailyExpenditures = computed(() => Array.from({ length: endOfMonth.value }, () => 500)) // Placeholder expenditure per day

const weekDays = ['Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab', 'Min']

const today = new Date()

const isToday = (day) => {
  const d = currentDate.value
  return (
    day === today.getDate() &&
    d.getMonth() === today.getMonth() &&
    d.getFullYear() === today.getFullYear()
  )
}

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
  <div class="mx-auto p-6 bg-white rounded-xl shadow-md">
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
    <table class="w-full border-collapse">
      <thead>
        <tr class="text-center text-sm font-medium text-gray-500">
          <th v-for="day in weekDays" :key="day" class="border-b border-gray-200 py-2">
            {{ day }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="week in weeksInMonth" :key="week.index" class="text-center">
          <td v-for="(day, dayIndex) in week.days" :key="dayIndex" :class="[
            'py-2 border border-gray-200 cursor-pointer',
            isToday(day)
              ? 'bg-blue-500 text-white font-semibold'
              : 'hover:bg-blue-100 text-gray-700'
          ]">
            <div v-if="day !== null" class="rounded-md border">
              {{ day }}
              <div class="flex justify-between mt-1">
                <div class="bg-green-500 text-white text-xs px-1 rounded flex-1">Budget: {{ dailyBudgets[day - 1]
                  }}</div>
                <div class="bg-red-500 text-white text-xs px-1 rounded flex-1">Pengeluaran: {{
                  dailyExpenditures[day - 1] }}
                </div>
              </div>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
