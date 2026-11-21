import { reactive } from 'vue';
import type { ProjectTimeline, TileServiceCategory } from '@/types/copely';

interface EstimateDraftState {
  service: TileServiceCategory;
  timeline: ProjectTimeline;
  highlightedFromHero: boolean;
}

const estimateDraft = reactive<EstimateDraftState>({
  service: 'porcelain-ceramic',
  timeline: 'this-week',
  highlightedFromHero: false,
});

export function useEstimateDraft() {
  function setDraftSelection(service: TileServiceCategory, timeline: ProjectTimeline): void {
    estimateDraft.service = service;
    estimateDraft.timeline = timeline;
    estimateDraft.highlightedFromHero = true;
  }

  function clearHighlight(): void {
    estimateDraft.highlightedFromHero = false;
  }

  return {
    estimateDraft,
    setDraftSelection,
    clearHighlight,
  };
}
