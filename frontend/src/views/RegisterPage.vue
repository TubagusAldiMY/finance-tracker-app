<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const handleRegister = () => {
  // Placeholder register logic
  console.log('Register with:', name.value, email.value, password.value)
  router.push('/dashboard')
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-100">
    <div class="bg-white p-8 rounded-2xl shadow-lg max-w-md w-full">
      <h1 class="text-2xl font-bold text-center mb-6">Register</h1>
      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">Name</label>
          <input v-model="name" type="text"
            class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="Your name">
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">Email</label>
          <input v-model="email" type="email"
            class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="email@example.com">
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">Password</label>
          <input v-model="password" type="password"
            class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="••••••••">
        </div>
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">Confirm Password</label>
          <input v-model="confirmPassword" type="password"
            class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="••••••••">
        </div>
        <button type="submit" class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
          Register
        </button>
      </form>
      <p class="text-center mt-4 text-slate-600">
        Already have an account?
        <router-link to="/login" class="text-blue-600 hover:underline">Login</router-link>
      </p>
    </div>
  </div>
</template>
