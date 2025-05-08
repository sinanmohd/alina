<script setup lang="ts">
const runtimeConfig = useRuntimeConfig();
const serverConfig = useServerConfig();
const serverUrl = useState('serverUrl', () => window.location.origin)

useSeoMeta({
  title: 'alina',
  ogTitle: 'alina',
  description: 'Your frenly neighbourhood file sharing website.',
  ogDescription: 'Your frenly neighbourhood file sharing website.',
})

await callOnce(async () => {
  if (runtimeConfig.public.serverUrl != '') {
    serverUrl.value = runtimeConfig.public.serverUrl;
  }

  const { data } = await useFetch(`${serverUrl.value}/_alina/config`)
  serverConfig.value = data.value as ServerConfig
});

console.log('ggg', useRuntimeConfig())
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
