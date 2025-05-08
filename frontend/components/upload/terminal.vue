<script setup lang="ts">
const appConfig = useAppConfig();
const { copyToClipboard } = useUtils();
const terminalExamble = `curl -F "file=@/path/to/your/file.jpg" ${appConfig.serverUrl}`;

function copyExamble(){
  copyToClipboard(terminalExamble);
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Terminal</CardTitle>
      <CardDescription>
        Upload files from the comfort of your terminal.
      </CardDescription>
    </CardHeader>

    <CardContent class="space-y-2">
    <code class="text-sm inline-flex bg-gray-50 rounded-lg p-2">
        <span class="flex gap-4">
            <span class="shrink-0 text-gray-500 my-auto">
                $
            </span>
            <span class="flex-1 wrap-anywhere my-auto">{{terminalExamble}}</span>
            <Button variant="outline" class="my-auto" :onclick="copyExamble">
              <Icon name="mdi:content-copy"/>
            </Button>
        </span>
    </code>
    </CardContent>
  </Card>
</template>
