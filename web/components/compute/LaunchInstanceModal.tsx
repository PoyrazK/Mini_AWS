
'use client';

import React, { useState } from 'react';
import { X } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import styles from './LaunchInstanceModal.module.css';

interface LaunchInstanceData {
  name: string;
  image: string;
  ports: string;
  vpc: string;
}

interface LaunchInstanceModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: LaunchInstanceData) => void;
}

export const LaunchInstanceModal: React.FC<LaunchInstanceModalProps> = ({ isOpen, onClose, onSubmit }) => {
  const [formData, setFormData] = useState<LaunchInstanceData>({
    name: '',
    image: 'ubuntu-22.04',
    ports: '80:80',
    vpc: 'default-vpc',
  });

  if (!isOpen) return null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
    onClose();
  };

  return (
    <div className={styles.overlay}>
      <div className={`${styles.modal} material-platter`}>
        <div className={styles.header}>
          <h2 className={styles.title}>Launch Instance</h2>
          <Button variant="ghost" onClick={onClose} className={styles.closeBtn}>
            <X size={18} />
          </Button>
        </div>
        
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label className={styles.label}>Name</label>
            <input 
              type="text" 
              className={styles.input} 
              placeholder="e.g. web-server-01"
              value={formData.name}
              onChange={(e) => setFormData({...formData, name: e.target.value})}
              autoFocus
            />
          </div>

          <div className={styles.row}>
            <div className={styles.field}>
              <label className={styles.label}>Image</label>
              <select 
                className={styles.select}
                value={formData.image}
                onChange={(e) => setFormData({...formData, image: e.target.value})}
              >
                <option value="ubuntu-22.04">Ubuntu 22.04 LTS</option>
                <option value="alpine-3.18">Alpine Linux 3.18</option>
                <option value="nginx-latest">Nginx (Latest)</option>
                <option value="postgres-15">PostgreSQL 15</option>
              </select>
            </div>

             <div className={styles.field}>
              <label className={styles.label}>VPC Network</label>
              <select 
                className={styles.select}
                value={formData.vpc}
                onChange={(e) => setFormData({...formData, vpc: e.target.value})}
              >
                <option value="default-vpc">default-vpc (172.31.0.0/16)</option>
                <option value="prod-net">prod-net (10.0.0.0/16)</option>
              </select>
            </div>
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Port Mapping (Host:Container)</label>
            <input 
              type="text" 
              className={styles.input} 
              placeholder="e.g. 8080:80, 443:443"
              value={formData.ports}
              onChange={(e) => setFormData({...formData, ports: e.target.value})}
            />
            <p className={styles.helpText}>Comma separated list of port mappings.</p>
          </div>

          <div className={styles.footer}>
            <Button type="button" variant="secondary" onClick={onClose}>Cancel</Button>
            <Button type="submit" variant="primary">Launch</Button>
          </div>
        </form>
      </div>
    </div>
  );
};
