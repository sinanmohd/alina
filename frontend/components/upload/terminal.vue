<script setup lang="ts">
const serverUrl = useState('serverUrl', () => window.location.origin)
const { copyToClipboard } = useUtils();
const terminalExamble = `curl -F "file=@/path/to/your/file.jpg" ${serverUrl.value}`;

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
      <span class="text-sm bg-gray-50 rounded-lg p-2 flex gap-4 justify-between w-full">
          <span class="text-gray-500 my-auto">
              $
          </span>

          <span class="flex-1 wrap-anywhere my-auto">{{terminalExamble}}</span>

          <Button variant="outline" class="my-auto" :onclick="copyExamble">
            <Icon name="lucide:copy"/>
          </Button>
      </span>
    </CardContent>
  </Card>
</template>
