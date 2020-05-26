import React from 'react';
import EntityPage from '../common/EntityPage';
import { api } from '../../common/api';

export const FundsPage = () => {
  return (
    <EntityPage
      entityName='Fund'
      entityPlural='Funds'
      blankEntity={{
        uuid: '',
        name: '',
        tickerSymbol: '',
        sharePrice: ''
      }}
      usesFilters={false}
      createEntity={api.createFund}
      getEntities={api.getFunds}
      updateEntity={api.updateFund}
      deleteEntity={api.deleteFund}
    />
  );
};