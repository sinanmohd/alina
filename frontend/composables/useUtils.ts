import JSZip from 'jszip'
import { toast } from 'vue-sonner';

export const useUtils = () => {
  return {
    formatData: (bytes: number, bits: boolean = false, decimals: number = 2): string => {
      if (bytes === 0) return '0 Bytes';

      const k = 1024;
      const dm = decimals < 0 ? 0 : decimals;
      let sizes;
      if (bits) {
        sizes = ['b', 'Kb', 'Mb', 'Gb', 'Tb', 'Pb', 'Eb', 'Zb', 'Yb'];
      } else {
        sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
      }

      const i = Math.floor(Math.log(bytes) / Math.log(k));

      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    },

    formatDuration: (speed: number, length: number): string => {
      if (speed <= 0 || length <= 0) {
        return 'Fast as possible';
      }

      const duration = length/speed;
      let durationFormatted: number
      let durationFormattedUnit: string

      if (duration < 60) {
        durationFormatted = duration
        durationFormattedUnit =  "second";
      } else if(duration < 60 * 60) {
        durationFormatted = (duration / 60)
        durationFormattedUnit =  "minute";
      } else if (duration < 60 * 60 * 24) {
        durationFormatted = (duration / (60 * 60))
        durationFormattedUnit =  "hour";
      } else {
        durationFormatted = (duration / (60 * 60 * 24))
        durationFormattedUnit =  "days";
      }

      return durationFormatted.toFixed(0) + " " + durationFormattedUnit + (duration < 2 ? "" : "s") + " left";
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
    },

    copyToClipboard: async (data: string) => {
      if (!navigator.clipboard) {
          toast("Clipboard Access Failed", {
            description: "For security reasons clipboard is disabled",
          });

          return;
      }

      await navigator.clipboard.writeText(data);
    },
  }
}
