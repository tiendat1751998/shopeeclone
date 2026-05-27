"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";

export function FlashSaleTimer() {
  const calculateTimeLeft = useCallback(() => {
    const endTime = new Date().getTime() + 2 * 60 * 60 * 1000;
    const now = new Date().getTime();
    const diff = endTime - now;

    if (diff <= 0) {
      return { hours: "00", minutes: "00", seconds: "00" };
    }

    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((diff % (1000 * 60)) / 1000);

    return {
      hours: String(hours).padStart(2, "0"),
      minutes: String(minutes).padStart(2, "0"),
      seconds: String(seconds).padStart(2, "0"),
    };
  }, []);

  const [timeLeft, setTimeLeft] = useState(calculateTimeLeft);

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft(calculateTimeLeft());
    }, 1000);
    return () => clearInterval(timer);
  }, [calculateTimeLeft]);

  return (
    <section className="mb-3">
      <div className="max-w-tiki mx-auto px-3">
        <div className="flash-sale px-4 py-2 flex items-center justify-between text-white rounded-lg">
          <div className="flex items-center gap-2">
            <span className="text-base">⚡</span>
            <span className="font-bold text-sm">FLASH SALE</span>
            <div className="flex items-center gap-1">
              <div className="bg-white/20 rounded px-1 py-0.5 text-xs font-mono font-bold">{timeLeft.hours}</div>
              <span>:</span>
              <div className="bg-white/20 rounded px-1 py-0.5 text-xs font-mono font-bold">{timeLeft.minutes}</div>
              <span>:</span>
              <div className="bg-white/20 rounded px-1 py-0.5 text-xs font-mono font-bold">{timeLeft.seconds}</div>
            </div>
          </div>
          <Link href="/promotions" className="text-xs font-medium text-white/90 hover:text-white">
            Xem tất cả →
          </Link>
        </div>
      </div>
    </section>
  );
}
