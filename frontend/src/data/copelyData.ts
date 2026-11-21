import type {
  ContentBundle,
  ProjectTimeline,
  TileServiceCategory,
} from '@/types/copely';
import { assetUrl } from '@/utils/assetUrl';

export interface SelectOption<T extends string = string> {
  value: T;
  label: string;
}

export const SERVICE_OPTIONS: SelectOption<TileServiceCategory>[] = [
  { value: 'porcelain-ceramic', label: 'Porcelain & Tile' },
  { value: 'natural-stone', label: 'Natural Stone & Marble' },
  { value: 'wet-rooms', label: 'Wet Rooms & Showers' },
  { value: 'kitchen-backsplashes', label: 'Kitchen Backsplashes' },
  { value: 'outdoor-patios', label: 'Outdoor Patios' },
];

export const TIMELINE_OPTIONS: SelectOption<ProjectTimeline>[] = [
  { value: 'this-week', label: 'This week' },
  { value: 'next-2-weeks', label: 'Next 2 weeks' },
  { value: 'planning-ahead', label: 'Planning ahead' },
];

export interface ComparisonCell {
  text: string;
  status: 'positive' | 'negative' | 'neutral';
}

export interface ComparisonRow {
  label: string;
  copely: ComparisonCell;
  diy: ComparisonCell;
  generalContractor: ComparisonCell;
  bigBox: ComparisonCell;
}

export const COMPARISON_ROWS: ComparisonRow[] = [
  {
    label: 'Cost',
    copely: { text: '$$', status: 'positive' },
    diy: { text: '$', status: 'neutral' },
    generalContractor: { text: '$$$', status: 'neutral' },
    bigBox: { text: '$$$$', status: 'negative' },
  },
  {
    label: 'Turnaround',
    copely: { text: '1–3 days', status: 'positive' },
    diy: { text: 'Weeks', status: 'neutral' },
    generalContractor: { text: 'Weeks', status: 'negative' },
    bigBox: { text: '4–5 days', status: 'neutral' },
  },
  {
    label: 'Quality',
    copely: { text: 'Guaranteed', status: 'positive' },
    diy: { text: 'Inconsistent', status: 'negative' },
    generalContractor: { text: 'Inconsistent', status: 'neutral' },
    bigBox: { text: 'Consistent', status: 'neutral' },
  },
  {
    label: 'Warranty',
    copely: { text: '25-year, transferable', status: 'positive' },
    diy: { text: 'None', status: 'negative' },
    generalContractor: { text: '1 year', status: 'neutral' },
    bigBox: { text: 'Limited', status: 'neutral' },
  },
  {
    label: 'Flexibility',
    copely: { text: 'Scheduled around you', status: 'positive' },
    diy: { text: 'Around their schedule', status: 'neutral' },
    generalContractor: { text: 'Their backlog', status: 'negative' },
    bigBox: { text: 'Fixed slots', status: 'neutral' },
  },
];

export interface ProcessStep {
  number: string;
  title: string;
  description: string;
}

export const PROCESS_STEPS: ProcessStep[] = [
  {
    number: '01',
    title: 'Free Inspection & Measure',
    description:
      'We laser-measure your rooms, check your subfloor, and tell you what prep it needs — then leave the decision to you.',
  },
  {
    number: '02',
    title: 'Clear Written Estimate',
    description:
      'You get a written estimate that lists every cost and your material choices. The price we quote is the price you pay.',
  },
  {
    number: '03',
    title: 'Pick Your Date',
    description:
      'We set a date that works for you and deliver your membranes and mortar the day before, so the crew starts on time.',
  },
  {
    number: '04',
    title: 'Installation',
    description:
      'Our licensed crew preps the substrate, waterproofs, and sets your new tile. Most rooms are done in two to four days.',
  },
  {
    number: '05',
    title: 'Final Walkthrough',
    description:
      'We grout, seal, polish every surface spotless, and walk the finished floor with you before we call it done.',
  },
];

export interface PolaroidPhoto {
  id: string;
  imageUrl: string;
  alt: string;
  transformClass: string;
}

export const FOOTER_POLAROIDS: PolaroidPhoto[] = [
  {
    id: 'pol-1',
    imageUrl: assetUrl('/images/polaroid-1.jpg'),
    alt: 'Sunlit stone bathroom remodel in Texas Hill Country',
    transformClass: 'rotate-[-8deg] translate-y-3 z-0',
  },
  {
    id: 'pol-2',
    imageUrl: assetUrl('/images/polaroid-2.jpg'),
    alt: 'Copely setter laying herringbone floor tile',
    transformClass: 'rotate-[-4deg] translate-y-1 z-10',
  },
  {
    id: 'pol-3',
    imageUrl: assetUrl('/images/polaroid-3.jpg'),
    alt: 'Craftsman troweling thin-set mortar in high-vis gear',
    transformClass: 'rotate-[0deg] -translate-y-2 z-20',
  },
  {
    id: 'pol-4',
    imageUrl: assetUrl('/images/polaroid-4.jpg'),
    alt: 'Illuminated modern home exterior and tiled patio at dusk',
    transformClass: 'rotate-[4deg] translate-y-1 z-10',
  },
  {
    id: 'pol-5',
    imageUrl: assetUrl('/images/polaroid-5.jpg'),
    alt: 'Copely crew van parked at a Hill Country residence',
    transformClass: 'rotate-[8deg] translate-y-3 z-0',
  },
];

export const FALLBACK_CONTENT: ContentBundle = {
  towns: [
    'Dripping Springs',
    'Austin',
    'Bee Cave',
    'Lakeway',
    'Wimberley',
    'Fredericksburg',
    'Boerne',
    'Marble Falls',
    'Johnson City',
    'San Marcos',
  ],
  services: [
    {
      id: 'svc-porcelain',
      slug: 'porcelain-ceramic',
      title: 'Porcelain & Tile',
      shortSubtitle: 'Large-format floors & walls',
      description:
        'Laser-leveled subfloor prep, uncoupling membranes, and large-format porcelain or ceramic tile set with razor-thin lippage-free grout joints.',
      imageUrl: assetUrl('/images/service-porcelain.jpg'),
      warrantyYears: 25,
      timelineLabel: '2–4 days typical',
      idealFor: 'Open-concept living rooms, kitchens, foyers, and high-traffic hallways.',
      specs: [
        'ANSI A108 substrate flattening to < 1/8" over 10 ft',
        'Crack-isolation uncoupling membrane over concrete slabs',
        '100% coverage back-buttering on tiles larger than 15"',
        'Stain-proof high-performance epoxy or urethane grout',
      ],
      processHighlights: [
        'Dust-extracted demo & slab grinding',
        'Self-leveling underlayment & movement joints',
        'Laser layout for balanced perimeter cuts',
      ],
    },
    {
      id: 'svc-wet-rooms',
      slug: 'wet-rooms',
      title: 'Wet Rooms',
      shortSubtitle: 'Seamless curbless showers & waterproofing',
      description:
        'Flood-tested bonded waterproof membranes, pre-sloped shower pans, linear drains, and custom tiled niches engineered to never leak.',
      imageUrl: assetUrl('/images/service-wetroom.jpg'),
      warrantyYears: 25,
      timelineLabel: '3–5 days typical',
      idealFor: 'Primary suite walk-in showers, curbless accessible wet rooms, and steam enclosures.',
      specs: [
        'Continuous sheet-membrane waterproofing wall-to-floor',
        '24-hour standing water flood test documented on-site',
        'Precision 1/4" per foot pan slope to linear or square drains',
        '100% silicone movement joints at all change-of-plane corners',
      ],
      processHighlights: [
        'Framing inspection & blocking for benches/glass',
        'Bonded pipe boots & mixing valve collar seals',
        'Mitered stone edges & recessed niche waterproofing',
      ],
    },
    {
      id: 'svc-backsplashes',
      slug: 'kitchen-backsplashes',
      title: 'Backsplashes',
      shortSubtitle: 'Custom kitchen mosaics, zellige & stone',
      description:
        'Artisan zellige, herringbone, and slab backsplashes finished around outlets, range hoods, and custom millwork with exact symmetry.',
      imageUrl: assetUrl('/images/service-backsplash.jpg'),
      warrantyYears: 25,
      timelineLabel: '1–2 days typical',
      idealFor: "Chef kitchens, butler's pantries, fireplace surrounds, and wet bars.",
      specs: [
        'Flush outlet box extensions & laser-centered focal layouts',
        'Color-matched anodized metal or mitered edge returns',
        'Non-sag polymer-modified mortar for handmade clay & glass',
        'Penetrating sealer applied prior to and after grouting',
      ],
      processHighlights: [
        'Full countertop & range protection before cut #1',
        'Dry-lay pattern blending for tonal variation',
        'Clean silicone perimeter caulk along countertops',
      ],
    },
    {
      id: 'svc-natural-stone',
      slug: 'natural-stone',
      title: 'Natural Stone',
      shortSubtitle: 'Marble, travertine, slate & limestone',
      description:
        'Honed limestone, veined Carrara marble, and Hill Country travertine installed with non-staining white thin-set and penetrating sealers.',
      imageUrl: assetUrl('/images/service-natural-stone.jpg'),
      warrantyYears: 25,
      timelineLabel: '3–5 days typical',
      idealFor: 'Architectural foyers, dining rooms, fireplace hearths, and primary baths.',
      specs: [
        'Crack-suppressing membrane beneath moisture-sensitive stone',
        'Non-staining white rapid-setting stone mortar',
        'Vein-matched dry layout approved with homeowner',
        'Breathable fluoropolymer penetrating stone sealer',
      ],
      processHighlights: [
        'Zero-deflection subfloor stiffening check',
        'Wet diamond polishing on exposed stone edges',
        'Care guide & neutral-pH maintenance kit included',
      ],
    },
    {
      id: 'svc-outdoor-patios',
      slug: 'outdoor-patios',
      title: 'Outdoor Patios',
      shortSubtitle: 'Frost-proof pavers, pool decks & steps',
      description:
        'Exterior porcelain pavers, natural stone terraces, and pool copings engineered with positive drainage and freeze-thaw resistant mortars.',
      imageUrl: assetUrl('/images/service-outdoor-patio.jpg'),
      warrantyYears: 25,
      timelineLabel: '3–6 days typical',
      idealFor: 'Covered porches, pool surrounds, outdoor kitchens, and hillside stone steps.',
      specs: [
        'Waterproofing & anti-efflorescence drainage mat',
        'Freeze-thaw rated exterior polymer mortar & epoxy grout',
        'Engineered perimeter expansion joints for Texas heat swings',
        'Slip-resistant R11 textured stone & exterior porcelain finishes',
      ],
      processHighlights: [
        'Slope verification away from home foundation',
        'Custom bullnose & mitered step treads',
        'Hydrophobic surface protection against oak tannins',
      ],
    },
  ],
  beforeAfterCases: [
    {
      id: 'ba-hero',
      title: 'Open-Concept Porcelain Slab Replacement',
      subtitle: 'Cracked 1990s builder ceramic replaced with 24×48 charcoal slate porcelain',
      location: 'Dripping Springs, TX',
      beforeImage: assetUrl('/images/ba-hero-before.jpg'),
      afterImage: assetUrl('/images/ba-hero-after.jpg'),
      featured: true,
    },
    {
      id: 'ba-bath',
      title: 'Primary Marble Wet Room',
      subtitle: 'Leaking fiberglass-pan shower rebuilt into curbless Carrara marble suite',
      location: 'Bee Cave, TX',
      beforeImage: assetUrl('/images/ba-bath-before.jpg'),
      afterImage: assetUrl('/images/ba-bath-after.jpg'),
      featured: false,
    },
    {
      id: 'ba-patio',
      title: 'Limestone & Slate Paver Terrace',
      subtitle: 'Weathered bare concrete deck transformed into custom Hill Country stone patio',
      location: 'Wimberley, TX',
      beforeImage: assetUrl('/images/ba-patio-before.jpg'),
      afterImage: assetUrl('/images/ba-patio-after.jpg'),
      featured: false,
    },
    {
      id: 'ba-kitchen',
      title: 'Artisan Glazed Herringbone Splashback',
      subtitle: 'Peeling dated backsplash replaced with handmade glazed herringbone tile',
      location: 'Lakeway, TX',
      beforeImage: assetUrl('/images/ba-kitchen-before.jpg'),
      afterImage: assetUrl('/images/ba-kitchen-after.jpg'),
      featured: false,
    },
  ],
  testimonials: [
    {
      id: 'test-daniel',
      author: 'Daniel Okafor',
      projectType: 'Tile repair',
      location: 'Houston, TX',
      avatarUrl: assetUrl('/images/avatar-1.jpg'),
      quote:
        'Two other companies told me I needed a full slab replacement. Copely found one bad expansion joint, fixed it in an afternoon, and charged me a few hundred dollars. They talked themselves out of a much bigger job.',
      heroImageUrl: assetUrl('/images/case-1.jpg'),
      metrics: [
        { label: 'Turnaround', value: '1 Afternoon' },
        { label: 'Transparency', value: '2 Quotes beaten' },
        { label: 'Protection', value: '25yr Warranty' },
      ],
    },
    {
      id: 'test-claire',
      author: 'Claire Morales',
      projectType: 'Whole-home porcelain',
      location: 'Dripping Springs, TX',
      avatarUrl: assetUrl('/images/avatar-2.jpg'),
      quote:
        'They leveled our slab the same week we called, and 1,200 sq ft of new large-format porcelain was down four days later. Every transition at the patio doors is dead flat.',
      heroImageUrl: assetUrl('/images/case-2.jpg'),
      metrics: [
        { label: 'Scope', value: '1,200 sq ft' },
        { label: 'Schedule', value: '4 Days' },
        { label: 'Protection', value: '25yr Warranty' },
      ],
    },
    {
      id: 'test-harold',
      author: 'Harold & Beth Vance',
      projectType: 'Curbless wet room',
      location: 'Fredericksburg, TX',
      avatarUrl: assetUrl('/images/avatar-3.jpg'),
      quote:
        'After a slow leak ruined our downstairs ceiling, Copely stripped the shower to studs, flood-tested the new pan for 24 hours, and set bookmatched stone in five days flat.',
      heroImageUrl: assetUrl('/images/case-3.jpg'),
      metrics: [
        { label: 'Flood Test', value: '24hr Verified' },
        { label: 'Duration', value: '5 Days' },
        { label: 'Protection', value: '25yr Warranty' },
      ],
    },
    {
      id: 'test-nadia',
      author: 'Nadia Al-Mansoor',
      projectType: 'Zellige kitchen backsplash',
      location: 'Austin, TX',
      avatarUrl: assetUrl('/images/avatar-4.jpg'),
      quote:
        "Our handmade Moroccan tile was tricky because every edge varies slightly. Sam's crew dry-laid the entire wall before touching mortar, and the finish around our brass hood is museum quality.",
      heroImageUrl: assetUrl('/images/case-4.jpg'),
      metrics: [
        { label: 'Duration', value: '1.5 Days' },
        { label: 'Edge Finish', value: 'Mitered Trim' },
        { label: 'Protection', value: '25yr Warranty' },
      ],
    },
    {
      id: 'test-graham',
      author: 'Graham Whitaker',
      projectType: 'Outdoor limestone patio',
      location: 'Wimberley, TX',
      avatarUrl: assetUrl('/images/avatar-5.jpg'),
      quote:
        'Texas summers and heavy storms wrecked our old flagstone mortar. Copely installed an uncoupling drainage mat and frost-proof stone pavers that still look brand new two winters later.',
      heroImageUrl: assetUrl('/images/case-5.jpg'),
      metrics: [
        { label: 'Area', value: '680 sq ft' },
        { label: 'System', value: 'Drainage Mat' },
        { label: 'Protection', value: '25yr Warranty' },
      ],
    },
  ],
  recentWork: [
    {
      id: 'rw-1',
      title: 'Notched Mortar Bed & Membrane Prep',
      location: 'Austin, TX',
      category: 'Substrate Prep',
      imageUrl: assetUrl('/images/bento-1.jpg'),
    },
    {
      id: 'rw-2',
      title: 'White Oak Herringbone Porcelain',
      location: 'Dripping Springs, TX',
      category: 'Porcelain & Tile',
      imageUrl: assetUrl('/images/bento-2.jpg'),
    },
    {
      id: 'rw-3',
      title: 'Custom Stone Bench & Niche Wet Room',
      location: 'Bee Cave, TX',
      category: 'Wet Rooms',
      imageUrl: assetUrl('/images/bento-3.jpg'),
    },
    {
      id: 'rw-4',
      title: 'Limestone Pool Coping & Terrace Steps',
      location: 'Wimberley, TX',
      category: 'Outdoor Patios',
      imageUrl: assetUrl('/images/bento-4.jpg'),
    },
    {
      id: 'rw-5',
      title: 'Precision Lippage-Free Leveling System',
      location: 'Lakeway, TX',
      category: 'Large-Format',
      imageUrl: assetUrl('/images/bento-5.jpg'),
    },
    {
      id: 'rw-6',
      title: 'Bookmatched Gallery Entryway Walls',
      location: 'Fredericksburg, TX',
      category: 'Natural Stone',
      imageUrl: assetUrl('/images/bento-6.jpg'),
    },
  ],
  homeFaqs: [
    {
      id: 'faq-1',
      category: 'estimates',
      question: 'How much does an estimate cost?',
      answer:
        'Nothing. We measure your space, take photos, and give you a written estimate — whether or not you hire us.',
    },
    {
      id: 'faq-2',
      category: 'during-job',
      question: 'How long does a full tile installation take?',
      answer:
        'Kitchen backsplashes typically take 1–2 days, most floor installations finish in 2–4 days, and complete waterproofed curbless wet rooms take 3–5 days including our 24-hour pan flood test.',
    },
    {
      id: 'faq-3',
      category: 'estimates',
      question: 'Do you handle water damage or insurance claims?',
      answer:
        'Yes. When a leaking shower pan or broken pipe damages subfloors or tile, we document moisture readings, supply line-item estimates compatible with major carriers, and bill insurance directly where approved.',
    },
    {
      id: 'faq-4',
      category: 'warranty',
      question: 'What does the warranty actually cover?',
      answer:
        'Our 25-year transferable warranty covers waterproofing integrity, mortar bond failure, substrate prep, lippage tolerances, and structural grout cracking — backed in writing for you and the next homeowner.',
    },
    {
      id: 'faq-5',
      category: 'waterproofing',
      question: 'Can you waterproof my shower before tiling?',
      answer:
        'Always. Tile and grout themselves are never waterproof. We install continuous bonded sheet waterproofing membranes, sealed valve collars, and sloped mortar beds before setting a single tile.',
    },
    {
      id: 'faq-6',
      category: 'during-job',
      question: 'What happens to all the old tile debris and dust?',
      answer:
        'We mask walkways with heavy floor protection, run HEPA dust-extraction shrouds on our grinders and demo hammers, and haul every scrap of old tile and mortar off-site daily.',
    },
  ],
  faqCategories: [
    {
      id: 'estimates',
      title: 'Estimates & pricing',
      items: [
        {
          id: 'faq-1',
          category: 'estimates',
          question: 'How much does an estimate cost?',
          answer:
            'Nothing. We measure your space, take photos, and give you a written estimate — whether or not you hire us.',
        },
        {
          id: 'faq-est-2',
          category: 'estimates',
          question: 'Can you supply the tile or work with tile I already bought?',
          answer:
            'Both. You can purchase tile through our Hill Country trade showrooms at contractor pricing, or supply your own material and we bring all membranes, mortar, trims, and grout.',
        },
        {
          id: 'faq-est-3',
          category: 'estimates',
          question: "Are there hidden fees if my slab isn't level?",
          answer:
            'No surprises. We laser-check floor flatness during your free estimate visit so self-leveling compound or joist stiffening is priced transparently upfront.',
        },
        {
          id: 'faq-3',
          category: 'estimates',
          question: 'Do you handle water damage or insurance claims?',
          answer:
            'Yes. When a leaking shower pan or broken pipe damages subfloors or tile, we document moisture readings, supply line-item estimates compatible with major carriers, and bill insurance directly where approved.',
        },
      ],
    },
    {
      id: 'during-job',
      title: 'During the job',
      items: [
        {
          id: 'faq-2',
          category: 'during-job',
          question: 'How long does a full tile installation take?',
          answer:
            'Kitchen backsplashes typically take 1–2 days, most floor installations finish in 2–4 days, and complete waterproofed curbless wet rooms take 3–5 days including our 24-hour pan flood test.',
        },
        {
          id: 'faq-6',
          category: 'during-job',
          question: 'What happens to all the old tile debris and dust?',
          answer:
            'We mask walkways with heavy floor protection, run HEPA dust-extraction shrouds on our grinders and demo hammers, and haul every scrap of old tile and mortar off-site daily.',
        },
        {
          id: 'faq-job-3',
          category: 'during-job',
          question: 'Can we stay in our home while you tile?',
          answer:
            'Yes. We set up zippered dust containment barriers, clean up tools at the end of every shift, and plan room sequencing so your kitchen or second bath stays usable.',
        },
        {
          id: 'faq-job-4',
          category: 'during-job',
          question: 'When can we walk on new floors or turn on the shower?',
          answer:
            'Light foot traffic is usually safe 12–18 hours after grouting. Wet rooms and steam showers are ready for full use 48 hours after final silicone curing.',
        },
      ],
    },
    {
      id: 'system',
      title: 'Your tile system, explained',
      items: [
        {
          id: 'faq-sys-1',
          category: 'system',
          question: 'Why do you use an uncoupling membrane over concrete?',
          answer:
            'Central Texas clay soil expands and contracts seasonally, creating hairline slab movement. An uncoupling membrane absorbs lateral shear so your porcelain or stone never cracks.',
        },
        {
          id: 'faq-sys-2',
          category: 'system',
          question: 'What makes large-format tile harder to install?',
          answer:
            'Tiles 15 inches or longer amplify minor subfloor dips and require 95%+ mortar coverage via back-buttering plus mechanical leveling clips to guarantee zero lippage.',
        },
        {
          id: 'faq-sys-3',
          category: 'system',
          question: 'Do you use cement grout or epoxy grout?',
          answer:
            'We recommend stain-proof high-performance urethane or epoxy grout in showers and kitchens so coffee, hard water, and soap never discolor grout joints.',
        },
      ],
    },
    {
      id: 'waterproofing',
      title: 'Waterproofing & prep',
      items: [
        {
          id: 'faq-5',
          category: 'waterproofing',
          question: 'Can you waterproof my shower before tiling?',
          answer:
            'Always. Tile and grout themselves are never waterproof. We install continuous bonded sheet waterproofing membranes, sealed valve collars, and sloped mortar beds before setting a single tile.',
        },
        {
          id: 'faq-wp-2',
          category: 'waterproofing',
          question: 'Do you flood-test every new shower pan?',
          answer:
            'Yes. Every custom wet room pan is plugged and filled for 24 hours to verify zero water drop before a single floor tile goes down.',
        },
        {
          id: 'faq-wp-3',
          category: 'waterproofing',
          question: 'How do you prevent mold in shower corners?',
          answer:
            'Instead of rigid grout at inside wall-to-floor corners—which always cracks—we bridge transitions with waterproof banding and finish with 100% mold-resistant silicone.',
        },
      ],
    },
    {
      id: 'warranty',
      title: 'Warranty & after',
      items: [
        {
          id: 'faq-4',
          category: 'warranty',
          question: 'What does the warranty actually cover?',
          answer:
            'Our 25-year transferable warranty covers waterproofing integrity, mortar bond failure, substrate prep, lippage tolerances, and structural grout cracking — backed in writing for you and the next homeowner.',
        },
        {
          id: 'faq-war-2',
          category: 'warranty',
          question: 'Is the 25-year warranty transferable if I sell my house?',
          answer:
            'Yes. Your signed warranty certificate transfers automatically to the next owner, adding verifiable value during home inspections.',
        },
        {
          id: 'faq-war-3',
          category: 'warranty',
          question: 'Do natural stone floors need periodic resealing?',
          answer:
            'Interior honed marble and limestone typically need a quick wipe-on penetrating sealer every 3–5 years; exterior stone patios benefit from sealing every 2–3 years.',
        },
      ],
    },
  ],
  googleReviews: [
    {
      id: 'rev-1',
      author: 'Max Rivera',
      initials: 'MR',
      rating: 5,
      dateLabel: '2 weeks ago',
      projectTag: 'Porcelain Living Floor',
      body: 'Copely ground down a high concrete hump in our living room before setting 24x48 porcelain. Zero hollow spots, zero lippage, and their crew swept the driveway before leaving each day.',
    },
    {
      id: 'rev-2',
      author: 'J. Whitfield',
      initials: 'JW',
      rating: 5,
      dateLabel: '1 month ago',
      projectTag: 'Curbless Wet Room',
      body: 'Night-and-day difference from our old contractor. They caught missing stud blocking, flood-tested the shower pan overnight, and showed us photos of every waterproofing layer before tiling.',
    },
    {
      id: 'rev-3',
      author: 'Gary Hutchins',
      initials: 'GH',
      rating: 5,
      dateLabel: '2 months ago',
      projectTag: 'Outdoor Pool Patio',
      body: 'Jordan walked us through drainage mats and exterior frost-proof pavers for our Wimberley deck. Finished right on the quoted price even after fixing a tricky threshold step.',
    },
    {
      id: 'rev-4',
      author: 'Janet Solomon',
      initials: 'JS',
      rating: 5,
      dateLabel: '3 months ago',
      projectTag: 'Kitchen Zellige Backsplash',
      body: "We worried our handmade zellige tiles would look wavy in a bad way. Sam's setters blended boxes carefully and mitered every outside corner. Looks like an architectural magazine.",
    },
    {
      id: 'rev-5',
      author: 'Trevor & Maya Lin',
      initials: 'TL',
      rating: 5,
      dateLabel: '4 months ago',
      projectTag: 'Whole-Home Stone',
      body: 'Three estimates in, Copely was the only crew who checked slab moisture and expansion joints with real instruments. Best home investment we have made.',
    },
    {
      id: 'rev-6',
      author: 'Elena Vance-Ramos',
      initials: 'EV',
      rating: 5,
      dateLabel: '5 months ago',
      projectTag: 'Shower Pan Repair',
      body: 'Honest, fast, and tidy. They saved our primary bath from a subfloor rot disaster and kept dust completely out of our bedroom closet.',
    },
  ],
  glossaryTerms: [
    {
      id: 'gl-1',
      term: 'Subfloor',
      definition:
        'The structural concrete slab or plywood decking beneath your tile that must be flat, rigid, and deflection-free before any setting begins.',
    },
    {
      id: 'gl-2',
      term: 'Uncoupling membrane',
      definition:
        'A flexible engineered mat placed between the subfloor and tile so tiny seasonal slab shifts never telegraph up as cracked tiles.',
    },
    {
      id: 'gl-3',
      term: 'Flashing / Profiles',
      definition:
        'Precision extruded aluminum or stainless edge trims that lock down exposed tile borders and protect corners from chipping.',
    },
    {
      id: 'gl-4',
      term: 'Drip edge / Threshold',
      definition:
        'A pitched stone sill or transition saddle engineered to shed water back into the shower and bridge differing floor heights cleanly.',
    },
    {
      id: 'gl-5',
      term: 'Shower pan slope',
      definition:
        'The exact 1/4-inch-per-foot pitch built into the mortar bed so every drop of shower water flows reliably toward the drain.',
    },
    {
      id: 'gl-6',
      term: 'Pipe boot / Mixing valve seal',
      definition:
        'Prefabricated waterproof collars sealed directly over plumbing penetrations so steam and pressurized spray stay behind the tile.',
    },
    {
      id: 'gl-7',
      term: 'Expansion joint',
      definition:
        'A flexible perimeter or field seam filled with matching silicone rather than rigid grout so tile expanses breathe safely with temperature swings.',
    },
    {
      id: 'gl-8',
      term: 'Back-buttering',
      definition:
        'Skimming a flat coat of thin-set mortar onto the back of every large tile before setting it into the combed floor bed for 100% bond coverage.',
    },
    {
      id: 'gl-9',
      term: 'Epoxy grout',
      definition:
        'A non-porous, chemical-resistant two-part grout that never needs sealing and resists hard water, wine, coffee, and mildew stains permanently.',
    },
  ],
  blogPosts: [
    {
      id: 'post-1',
      slug: 'tile-installation-cost-texas-hill-country',
      category: 'PRICING',
      publishedAt: 'JUN 18, 2026',
      readTime: '6 MIN',
      title: 'What does new tile installation cost in the Hill Country?',
      excerpt:
        'A line-by-line look at subfloor flattening, waterproofing membranes, mortar specs, and labor rates across Austin, Dripping Springs, and Fredericksburg.',
      authorName: 'Jordan Pike',
      authorRole: 'Lead Estimator',
      authorAvatar: assetUrl('/images/author-jordan.jpg'),
      coverImage: assetUrl('/images/blog-featured.jpg'),
      featured: true,
      sections: [
        {
          heading: 'Why two tile quotes rarely price the same system',
          paragraphs: [
            'When homeowners collect three estimates for a primary bath or 1,000-square-foot living room floor, quotes can swing by 40% before anyone picks out a tile. Almost every gap hides beneath the surface: whether the installer includes slab grinding, crack-isolation membranes, and flood-tested waterproofing.',
            'At Copely, we itemize subfloor prep, uncoupling membranes, trims, and grout separately from tile selection so you know exactly what protects your floor for the next 25 years.',
          ],
        },
        {
          heading: 'Typical 2026 installed ranges by project type',
          paragraphs: [
            'Kitchen backsplashes in Austin and Dripping Springs generally fall between $2,200 and $4,800 depending on layout geometry (straight stack vs. herringbone) and outlet cuts. Large-format porcelain floor installations typically range from $14 to $24 per square foot installed including uncoupling membrane.',
            'Complete curbless wet room showers—including demolition, framing corrections, bonded sheet waterproofing, 24-hour flood testing, and custom stone benches—typically land between $7,500 and $14,000.',
          ],
        },
        {
          heading: 'Three questions to verify before signing any bid',
          paragraphs: [
            'Always ask whether your installer checks concrete flatness with a 10-foot straightedge, installs soft silicone movement joints at walls instead of hard grout, and backs waterproofing in writing.',
          ],
        },
      ],
    },
    {
      id: 'post-2',
      slug: '5-signs-shower-pan-hidden-water-damage',
      category: 'WATERPROOFING',
      publishedAt: 'MAY 27, 2026',
      readTime: '5 MIN',
      title: '5 signs your shower pan has hidden water damage',
      excerpt:
        'Musty baseboard odors, hollow tiles near the drain, and white efflorescence along curb joints are early warnings before subfloor rot starts.',
      authorName: 'Sam Calloway',
      authorRole: 'Founder & Master Setter',
      authorAvatar: assetUrl('/images/author-sam.jpg'),
      coverImage: assetUrl('/images/blog-waterproofing.jpg'),
      featured: false,
      sections: [
        {
          heading: 'Water travels quietly long before a ceiling stain appears',
          paragraphs: [
            'Most failed showers do not burst overnight. Capillary moisture seeps past cracked corner grout or unsealed curb fasteners into framing months before drywall shows discoloration.',
          ],
        },
        {
          heading: 'What to inspect around your shower today',
          paragraphs: [
            'Look for swollen baseboards outside the glass door, persistent damp grout lines 24 hours after showering, chalky mineral crust (efflorescence) along the curb, hollow-sounding floor tiles when tapped, or soft drywall beside the valve trim.',
          ],
        },
      ],
    },
    {
      id: 'post-3',
      slug: 'porcelain-vs-natural-stone-texas-homes',
      category: 'MATERIALS',
      publishedAt: 'MAY 8, 2026',
      readTime: '7 MIN',
      title: 'Porcelain vs. natural stone in Texas homes',
      excerpt:
        'How porosity, UV exposure, Hill Country hard water, and thermal expansion affect limestone, travertine, marble, and kiln-fired porcelain.',
      authorName: 'Erin Vasquez',
      authorRole: 'Co-Founder & Head of Operations',
      authorAvatar: assetUrl('/images/author-erin.jpg'),
      coverImage: assetUrl('/images/blog-materials.jpg'),
      featured: false,
      sections: [
        {
          heading: 'Porcelain: Zero-sealing resilience for busy spaces',
          paragraphs: [
            'Modern rectified Italian and Spanish porcelain replicates veined marble and warm limestone with <0.5% water absorption. It shrugs off dog claws, patio UV rays, and acidic kitchen spills with minimal maintenance.',
          ],
        },
        {
          heading: 'Natural Stone: Warmth and timeless patina',
          paragraphs: [
            'Nothing matches the organic depth of honed travertine or Texas limestone underfoot. Paired with a crack-suppression membrane and breathable penetrating sealer, natural stone wears gracefully for generations.',
          ],
        },
      ],
    },
    {
      id: 'post-4',
      slug: 'regrout-or-replace-how-to-decide',
      category: 'REPAIR',
      publishedAt: 'APR 14, 2026',
      readTime: '5 MIN',
      title: 'Regrout or replace? How to decide',
      excerpt:
        'When diamond-blade regrouting and movement-joint restoration save thousands—and when cracked tiles point to a deeper substrate failure.',
      authorName: 'Sam Calloway',
      authorRole: 'Founder & Master Setter',
      authorAvatar: assetUrl('/images/author-sam.jpg'),
      coverImage: assetUrl('/images/blog-repair.jpg'),
      featured: false,
      sections: [
        {
          heading: 'When targeted restoration beats tear-out',
          paragraphs: [
            'If your tiles sound solid when tapped and only surface cement grout has stained or crumbled—especially at perimeter expansion joints—removing old grout to full depth and refilling with urethane or epoxy grout restores a floor in a single day.',
          ],
        },
        {
          heading: 'When replacement is the only honest fix',
          paragraphs: [
            'If multiple tiles tent upward, rock underfoot, or crack in a continuous jagged line across a room, the underlying mortar bond or slab joint has failed. Regrouting alone will crack again within months.',
          ],
        },
      ],
    },
  ],
};
