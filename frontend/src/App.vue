<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { EventsOn } from '../wailsjs/runtime/runtime'

const router = useRouter()

onMounted(() => {
  EventsOn('navigate', (route: string) => {
    router.push(route)
  })
})
</script>
