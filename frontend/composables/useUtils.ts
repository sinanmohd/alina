import JSZip from 'jszip'

export const useUtils = () => {
  return {
    formatBytes: (bytes: number, decimals: number = 2): string => {
      if (bytes === 0) return '0 Bytes';

      const k = 1024;
      const dm = decimals < 0 ? 0 : decimals;
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

      const i = Math.floor(Math.log(bytes) / Math.log(k));

      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    },

    filesZip: async (files: FileList): Promise<Blob> => {
      const zip = new JSZip();

      for (const file of files) {
        zip.file(file.name, await file.arrayBuffer());
      }
      const blob = await zip.generateAsync({type: "blob"});

      return blob;
    },

    chunksFromBlob: (blob: Blob, chunkSize: number): Blob[] => {
      const chunks = Array<Blob>();

      for (let offset = 0; offset < blob.size; offset += chunkSize) {
        const end = offset + chunkSize > blob.size ? blob.size : offset + chunkSize;
        const chunk = blob.slice(offset, end);

        chunks.push(chunk);
      }

      return chunks
    }
  }
}
