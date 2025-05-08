<script setup lang="ts">
const appConfig = useAppConfig();
const serverConfig = useServerConfig();

useSeoMeta({
  title: 'alina',
  ogTitle: 'alina',
  description: 'Your frenly neighbourhood file sharing website.',
  ogDescription: 'Your frenly neighbourhood file sharing website.',
})

await callOnce(async () => {
  const { data } = await useFetch(`${appConfig.serverUrl}/_alina/config`)
  serverConfig.value = data.value as ServerConfig
});
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
