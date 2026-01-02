<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// Props for customization
const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  autoPlayInterval: {
    type: Number,
    default: 5000
  },
  subjects: {
    type: Array,
    required: true
  }
})

const currentSlide = ref(0)
let slideInterval = null

const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % props.subjects.length
}

const prevSlide = () => {
  currentSlide.value = currentSlide.value === 0
    ? props.subjects.length - 1
    : currentSlide.value - 1
}

const goToSlide = (index) => {
  currentSlide.value = index
}

onMounted(() => {
  if (props.autoPlayInterval > 0) {
    slideInterval = setInterval(nextSlide, props.autoPlayInterval)
  }
})

onUnmounted(() => {
  if (slideInterval) clearInterval(slideInterval)
})
</script>

<template>
  <section class="py-20 px-4 bg-white">
    <div class="max-w-5xl mx-auto">
      <h2 class="text-3xl font-bold text-center text-slate-800 mb-4">
        {{ title }}
      </h2>
      <p class="text-center text-slate-600 mb-12">
        {{ subtitle }}
      </p>

      <div class="relative">
        <div class="overflow-hidden rounded-2xl shadow-2xl bg-slate-100">
          <div class="flex transition-transform duration-500 ease-in-out"
            :style="{ transform: `translateX(-${currentSlide * 100}%)` }">
            <div v-for="obj in subjects" :key="obj.id" class="w-full flex-shrink-0">
              <div class="aspect-video flex items-center justify-center bg-slate-200 relative">
                <div class="text-center p-8">
                  <div class="w-24 h-24 mx-auto mb-4 bg-slate-300 rounded-full flex items-center justify-center">
                    <span class="text-4xl">📱</span>
                  </div>
                  <h3 class="text-xl font-semibold text-slate-600 mb-2">
                    {{ obj.title }}
                  </h3>
                  <p class="text-slate-500">
                    {{ obj.description }}
                  </p>
                  <div v-if="obj.image" class="aspect-video bg-slate-200 relative mt-4 rounded-lg overflow-hidden">
                    <img :src="obj.image" :alt="obj.title" class="w-full h-full object-contain" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Navigation buttons -->
        <button @click="prevSlide"
          class="absolute left-4 top-1/2 -translate-y-1/2 w-12 h-12 bg-white rounded-full shadow-lg flex items-center justify-center hover:bg-slate-50 transition z-10">
          <span class="text-xl">←</span>
        </button>
        <button @click="nextSlide"
          class="absolute right-4 top-1/2 -translate-y-1/2 w-12 h-12 bg-white rounded-full shadow-lg flex items-center justify-center hover:bg-slate-50 transition z-10">
          <span class="text-xl">→</span>
        </button>

        <!-- Dots indicator -->
        <div class="flex justify-center gap-2 mt-6">
          <button v-for="(screenshot, index) in subjects" :key="screenshot.id" @click="goToSlide(index)"
            class="w-3 h-3 rounded-full transition-all duration-300"
            :class="currentSlide === index ? 'bg-blue-600 w-8' : 'bg-slate-300'"></button>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
