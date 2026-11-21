import { ref } from 'vue';
import { FALLBACK_CONTENT } from '@/data/copelyData';
import { fetchContent } from '@/services/api';
import type { ContentBundle } from '@/types/copely';

const contentState = ref<ContentBundle>(FALLBACK_CONTENT);
const isLoaded = ref(false);
let inFlightPromise: Promise<ContentBundle> | null = null;

export function useCopelyContent() {
  async function loadContent(): Promise<ContentBundle> {
    if (isLoaded.value) {
      return contentState.value;
    }
    if (inFlightPromise) {
      return inFlightPromise;
    }

    inFlightPromise = fetchContent()
      .then((bundle) => {
        contentState.value = bundle;
        isLoaded.value = true;
        return bundle;
      })
      .catch(() => {
        // Retain typed fallback content seamlessly if API is unreachable during static preview.
        return contentState.value;
      })
      .finally(() => {
        inFlightPromise = null;
      });

    return inFlightPromise;
  }

  return {
    content: contentState,
    isLoaded,
    loadContent,
  };
}
