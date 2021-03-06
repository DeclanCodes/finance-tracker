import React, { useState } from 'react';
import { NavItem, NavItemProps } from '../NavItem/NavItem';
import './NavDropdown.scss';

interface NavDropdownProps {
  title: string,
  navItems: NavItemProps[]
}

export const NavDropdown = ({
  title,
  navItems
}: NavDropdownProps) => {
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const onClick = (_: any) => {
    setIsOpen(!isOpen);
  };

  return (
    <li
      className='nav-dropdown-container'
      tabIndex={0}
      onClick={onClick}
    >
      <div className='dropdown'>
        <div className='title'>{title}</div>
        <div className='triangle'/>
      </div>
      {isOpen &&
        <ul className='nav-items'>
          {navItems.map(ni => (
            <NavItem
              key={`nav-dropdown-item-${ni.title}`}
              to={ni.to}
              title={ni.title}
            />
          ))}
        </ul>
      }
    </li>
  );
};