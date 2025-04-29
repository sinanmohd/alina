<template>
  <Tabs default-value="files" class="w-full sm:w-3xl mx-auto p-4">

    <TabsList class="grid w-full grid-cols-3">
      <TabsTrigger value="files">
        <Icon name="mdi:files"/>
        Files
      </TabsTrigger>
      <TabsTrigger value="text">
        <Icon name="mdi:text"/>
        Text
      </TabsTrigger>
      <TabsTrigger value="terminal">
        <Icon name="mdi:terminal"/>
        Terminal
      </TabsTrigger>
    </TabsList>

    <TabsContent value="files">
      <UploadFiles />
    </TabsContent>

    <TabsContent value="text">
      <UploadText />
    </TabsContent>

    <TabsContent value="terminal">
      <UploadTerminal />
    </TabsContent>
  </Tabs>
</template>
