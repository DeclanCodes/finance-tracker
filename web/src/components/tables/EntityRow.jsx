import React, { useState } from 'react';
import { EntityForm } from '../forms/EntityForm';
import { ModifyRowPanel } from './ModifyRowPanel';
import { helpers } from '../../common/helpers';

export const EntityRow = ({
  entityName,
  entity,
  getInitialValues,
  options,
  doExtraModifications,
  handleUpdate,
  handleDelete
}) => {
  const [isEditing, setIsEditing] = useState(false);
  return (
    <tr>
      {entity.hasOwnProperty('name') && <td>{entity.name}</td>}
      {entity.hasOwnProperty('account') && <td>{entity.account.name}</td>}
      {entity.hasOwnProperty('category') && <td>{entity.category.name}</td>}
      {entity.hasOwnProperty('fund') && <td>{entity.fund.tickerSymbol}</td>}
      {entity.hasOwnProperty('description') && <td>{entity.description}</td>}
      {entity.hasOwnProperty('tickerSymbol') && <td>{entity.tickerSymbol}</td>}
      {entity.hasOwnProperty('date') && <td>{helpers.displayDate(entity.date)}</td>}
      {entity.hasOwnProperty('amount') && <td>{helpers.displayCurrency(entity.amount)}</td>}
      {entity.hasOwnProperty('sharePrice') && <td>{helpers.displayCurrency(entity.sharePrice)}</td>}
      {entity.hasOwnProperty('shares') && <td>{helpers.displayDecimals(entity.shares, 3)}</td>}
      {entity.hasOwnProperty('value') && <td>{helpers.displayCurrency(entity.value)}</td>}
      {entity.hasOwnProperty('expenseRatio') && <td>{`${helpers.displayPercentage(entity.expenseRatio, 3)}%`}</td>}
      {entity.hasOwnProperty('effectiveExpense') && <td>{helpers.displayCurrency(entity.effectiveExpense)}</td>}
      <td>
        {isEditing ? (
          <EntityForm
            entityName={entityName}
            entity={entity}
            getInitialValues={getInitialValues}
            options={options}
            doExtraModifications={doExtraModifications}
            doSubmit={handleUpdate}
            doFinalState={() => setIsEditing(false)}
          />
        ) : (
          <ModifyRowPanel
            handleEdit={() => setIsEditing(true)}
            handleDelete={() => handleDelete(entity.uuid)}
          />
        )}
      </td>
    </tr>
  );
};