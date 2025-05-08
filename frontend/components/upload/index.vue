<template>
  <Tabs default-value="files" class="w-full sm:w-3xl mx-auto p-4">

    <TabsList class="grid w-full grid-cols-2">
      <TabsTrigger value="files">
        <Icon name="mdi:files"/>
        Files
      </TabsTrigger>
      <TabsTrigger value="text">
        <Icon name="mdi:text"/>
        Text
      </TabsTrigger>
    </TabsList>

    <TabsContent value="files">
      <UploadFiles />
    </TabsContent>

    <TabsContent value="text">
      <UploadText />
    </TabsContent>
  </Tabs>
</template>
