import { useState, useEffect } from 'react';
import styles from '../styles/modal.module.css';

interface FilterModalProps {
  isOpen: boolean;
  onClose: () => void;
  onApply: (values: string[]) => void;
  title: string;
  options: { value: string; label: string }[];
  selectedValues: string[];
}

export const FilterModal = ({ isOpen, onClose, onApply, title, options, selectedValues }: FilterModalProps) => {
  const [tempSelected, setTempSelected] = useState<string[]>(selectedValues);

  useEffect(() => {
    setTempSelected(selectedValues);
  }, [selectedValues, isOpen]);

  const toggleOption = (value: string) => {
    setTempSelected((prev) =>
      prev.includes(value) ? prev.filter((v) => v !== value) : [...prev, value]
    );
  };

  const handleApply = () => {
    onApply(tempSelected);
    onClose();
  };

  const handleClear = () => {
    setTempSelected([]);
  };

  if (!isOpen) return null;

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.header}>
          <h3 className={styles.title}>{title}</h3>
        </div>

        <div className={styles.content}>
          {options.map((option) => (
            <label key={option.value} className={styles.option}>
              <input
                type="checkbox"
                checked={tempSelected.includes(option.value)}
                onChange={() => toggleOption(option.value)}
                className={styles.checkbox}
              />
              <span className={styles.label}>{option.label}</span>
            </label>
          ))}
        </div>

        <div className={styles.footer}>
          <button className={styles.clearBtn} onClick={handleClear}>
            Effacer
          </button>
          <button className={styles.applyBtn} onClick={handleApply}>
            Appliquer
          </button>
        </div>
      </div>
    </div>
  );
};
