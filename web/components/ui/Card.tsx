
import React from 'react';
import styles from './Card.module.css';

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  title?: string;
}

export const Card: React.FC<CardProps> = ({ children, className, title, ...props }) => {
  return (
    <div className={`${styles.card} material-platter ${className || ''}`} {...props}>
      {title && <div className={styles.header}>{title}</div>}
      <div className={styles.content}>
        {children}
      </div>
    </div>
  );
};
