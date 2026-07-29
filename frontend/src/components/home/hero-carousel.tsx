"use client";

import { useState, useEffect, useCallback } from "react";
import Image from "next/image";
import { Link } from "@/i18n/navigation";
import { Button } from "@/components/ui/button";
import { IconChevronLeft, IconChevronRight, IconArrowRight } from "@tabler/icons-react";

export interface HeroSlide {
  id: string;
  brand: string;
  title: string;
  subtitle: string;
  imageUrl: string;
  ctaText: string;
  ctaHref: string;
}

const DEFAULT_SLIDES: HeroSlide[] = [
  {
    id: "slide-1",
    brand: "STAND OIL & TRENDING",
    title: "SNEAKER COLLECTION",
    subtitle: "Compare prices from 40+ local & international brands, all in one place.",
    imageUrl: "/hero/slide-1.png",
    ctaText: "Find Sneakers",
    ctaHref: "/cari",
  },
  {
    id: "slide-2",
    brand: "NIKE & ADIDAS TECH",
    title: "PERFORMANCE RUNNING",
    subtitle: "Find the best price on top running shoes from verified sellers.",
    imageUrl: "/hero/slide-2.png",
    ctaText: "Shop Performance",
    ctaHref: "/cari?category=running",
  },
  {
    id: "slide-3",
    brand: "VINTAGE & STREETWEAR",
    title: "RETRO CLASSICS",
    subtitle: "Iconic classic releases from Compass, Ventela, to Vans & Converse.",
    imageUrl: "/hero/slide-3.png",
    ctaText: "Explore Brands",
    ctaHref: "/brand",
  },
];

interface HeroCarouselProps {
  slides?: HeroSlide[];
}

export function HeroCarousel({ slides = DEFAULT_SLIDES }: HeroCarouselProps) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isPaused, setIsPaused] = useState(false);

  const nextSlide = useCallback(() => {
    setCurrentIndex((prevIndex) => (prevIndex + 1) % slides.length);
  }, [slides.length]);

  const prevSlide = useCallback(() => {
    setCurrentIndex((prevIndex) => (prevIndex - 1 + slides.length) % slides.length);
  }, [slides.length]);

  useEffect(() => {
    if (isPaused) return;
    const interval = setInterval(() => {
      nextSlide();
    }, 5000);
    return () => clearInterval(interval);
  }, [nextSlide, isPaused]);

  const currentSlide = slides[currentIndex];

  return (
    <section className="mx-auto max-w-[1400px] px-4 py-6 md:px-6 md:py-8">
      {/* Full-Image Hero Banner Container */}
      <div 
        className="group relative h-[420px] md:h-[500px] w-full overflow-hidden rounded-2xl md:rounded-3xl border border-border/40 shadow-xl transition-all duration-500"
        onMouseEnter={() => setIsPaused(true)}
        onMouseLeave={() => setIsPaused(false)}
      >
        {/* Background Full Image */}
        <Image
          key={currentSlide.id}
          src={currentSlide.imageUrl}
          alt={currentSlide.title}
          fill
          priority
          sizes="100vw"
          className="object-cover transition-transform duration-700 group-hover:scale-105"
        />

        {/* Gradient Overlay for Text Readability */}
        <div className="absolute inset-0 bg-gradient-to-r from-black/85 via-black/50 to-transparent" />

        {/* Seamless Glassmorphism Left Arrow Button */}
        <button
          onClick={prevSlide}
          aria-label="Previous Slide"
          className="absolute left-4 top-1/2 z-20 -translate-y-1/2 -translate-x-4 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 flex h-11 w-11 items-center justify-center rounded-full border border-white/25 bg-white/15 text-white backdrop-blur-md transition-all duration-300 hover:scale-110 hover:bg-white hover:text-black shadow-lg md:left-6 md:h-12 md:w-12"
        >
          <IconChevronLeft className="h-6 w-6" />
        </button>

        {/* Seamless Glassmorphism Right Arrow Button */}
        <button
          onClick={nextSlide}
          aria-label="Next Slide"
          className="absolute right-4 top-1/2 z-20 -translate-y-1/2 translate-x-4 opacity-0 group-hover:translate-x-0 group-hover:opacity-100 flex h-11 w-11 items-center justify-center rounded-full border border-white/25 bg-white/15 text-white backdrop-blur-md transition-all duration-300 hover:scale-110 hover:bg-white hover:text-black shadow-lg md:right-6 md:h-12 md:w-12"
        >
          <IconChevronRight className="h-6 w-6" />
        </button>

        {/* Overlay Content */}
        <div className="relative z-10 flex h-full flex-col justify-center space-y-4 px-12 py-10 text-white sm:px-16 md:max-w-[65%] md:px-20 md:py-14">
          <div>
            <h1 className="font-display text-3xl font-bold uppercase leading-[0.95] tracking-tight text-white sm:text-5xl md:text-6xl lg:text-7xl">
              {currentSlide.title}
            </h1>
          </div>

          <p className="max-w-xl text-sm text-white/80 md:text-base">
            {currentSlide.subtitle}
          </p>

          <div className="pt-2">
            <Link href={currentSlide.ctaHref}>
              <Button size="default" className="gap-2 bg-white text-black hover:bg-white/90 font-medium px-6">
                {currentSlide.ctaText}
                <IconArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>

        {/* Slide Indicators Dots (Bottom Center) */}
        <div className="absolute bottom-4 left-1/2 z-20 flex -translate-x-1/2 gap-2">
          {slides.map((_, i) => (
            <button
              key={i}
              onClick={() => setCurrentIndex(i)}
              aria-label={`Go to slide ${i + 1}`}
              className={`h-1.5 rounded-sm transition-all duration-300 ${
                i === currentIndex ? "w-6 bg-white" : "w-2 bg-white/40 hover:bg-white/70"
              }`}
            />
          ))}
        </div>
      </div>
    </section>
  );
}
