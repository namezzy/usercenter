import { useState, useEffect } from 'react';

export const useMediaQuery = (query: string): boolean => {
  const [matches, setMatches] = useState(false);

  useEffect(() => {
    const media = window.matchMedia(query);
    
    const updateMatch = () => {
      setMatches(media.matches);
    };

    // 初始检查
    updateMatch();

    // 监听变化
    media.addEventListener('change', updateMatch);

    // 清理函数
    return () => {
      media.removeEventListener('change', updateMatch);
    };
  }, [query]);

  return matches;
};
