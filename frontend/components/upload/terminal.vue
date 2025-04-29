<template>
  <Card>
    <CardHeader>
      <CardTitle>Terminal</CardTitle>
      <CardDescription>
        Upload files from the comfert of your homly terminal.
      </CardDescription>
    </CardHeader>

    <CardContent class="space-y-2">
    <code class="text-sm inline-flex bg-gray-50 rounded-lg p-2">
        <span class="flex gap-4">
            <span class="shrink-0 text-gray-500 my-auto">
                $
            </span>

            <span class="flex-1 wrap-anywhere my-auto">
              curl -F "file=@/path/to/your/file.jpg" https://share.sinanmohd.com/
            </span>
            <Button variant="outline" class="my-auto">
              <Icon name="mdi:content-copy"/>
            </Button>
        </span>
    </code>
    </CardContent>
  </Card>
</template>
